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
	cpuUsageTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage_total",
		Help: "Total CPU usage percentage",
	})

	cpuUsagePerCore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_per_core",
			Help: "CPU usage percentage per core",
		},
		[]string{"core"},
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
	prometheus.MustRegister(cpuUsageTotal)
	prometheus.MustRegister(cpuUsagePerCore)
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
			var total float64
			for i, percentage := range percentages {
				cpuUsagePerCore.WithLabelValues(fmt.Sprintf("%d", i)).Set(percentage)
				total += percentage
			}
			cpuUsageTotal.Set(total / float64(len(percentages)))
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

	// Setup HTTP server
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting server on :2112")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
