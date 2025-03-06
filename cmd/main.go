package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var (
	// CPU metrics
	cpuUsagePerCore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_per_core",
			Help: "CPU usage percentage per core",
		},
		[]string{"core"},
	)

	cpuUsageSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "cpu_usage_summary",
			Help:       "Summary of overall CPU usage percentage over time",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // 记录中位数、90分位和99分位数
			MaxAge:     time.Minute * 5,                                        // 5分钟的时间窗口
		},
	)

	// Memory metrics
	memoryTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_total_bytes",
		Help: "Total memory in bytes",
	})

	memoryUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_used_bytes",
		Help: "Used memory in bytes",
	})

	memoryFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_free_bytes",
		Help: "Free memory in bytes",
	})

	memoryAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_available_bytes",
		Help: "Available memory in bytes",
	})
)

func init() {
	// Register all metrics
	prometheus.MustRegister(cpuUsagePerCore)
	prometheus.MustRegister(cpuUsageSummary)
	prometheus.MustRegister(memoryTotal)
	prometheus.MustRegister(memoryUsed)
	prometheus.MustRegister(memoryFree)
	prometheus.MustRegister(memoryAvailable)
}

func collectMetrics() {
	for {
		// Collect CPU usage
		percentages, err := cpu.Percent(time.Second, true)
		if err != nil {
			log.Printf("Error collecting CPU metrics: %v", err)
		} else {
			// Record per-core metrics
			for i, percentage := range percentages {
				core := fmt.Sprintf("%d", i)
				cpuUsagePerCore.WithLabelValues(core).Set(percentage)
			}

			// Calculate and record average CPU usage for summary
			var total float64
			for _, percentage := range percentages {
				total += percentage
			}
			average := total / float64(len(percentages))
			cpuUsageSummary.Observe(average)
		}

		// Collect memory information
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Error collecting memory metrics: %v", err)
		} else {
			memoryTotal.Set(float64(memInfo.Total))
			memoryUsed.Set(float64(memInfo.Used))
			memoryFree.Set(float64(memInfo.Free))
			memoryAvailable.Set(float64(memInfo.Available))
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	// Start metrics collection
	go collectMetrics()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting server on :2112")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
