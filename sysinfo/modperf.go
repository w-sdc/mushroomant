package sysinfo

import (
	"context"
	"time"
)

type CPUStat struct {
	Total float32   `json:"total"`
	Core  []float32 `json:"core"`
}

type MemStat struct {
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
	Available uint64 `json:"available"`
}

type NetStat struct {
	BytesSend   uint64 `json:"bytes_send"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSend uint64 `json:"packets_send"`
	PacketsRecv uint64 `json:"packets_recv"`
}

type DiskUsage struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

type ContainerStat struct {
	CPU     float32 `json:"cpu"`
	MemUsed uint64  `json:"mem_used"`
}

type ContianerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Runing bool   `json:"runing"`
}

// PerfStat is a struct that contains all the performance data
// To be an event, it might only contain the data that is changed
type PerfStat struct {
	CPU       *CPUStat                 `json:"cpu,omitempty"`
	Mem       *MemStat                 `json:"mem,omitempty"`
	NetIOPSec map[string]NetStat       `json:"net_io,omitempty"`
	DiskUsage map[string]DiskUsage     `json:"disk_usage,omitempty"`
	CStat     map[string]ContainerStat `json:"c_stat,omitempty"`
	CEvent    []ContianerInfo          `json:"c_event,omitempty"`
}

type PerfTimeline struct {
	Interval int64           `json:"interval"`
	Stats    []PerfStat      `json:"stats"`
	CInfo    []ContianerInfo `json:"c_event"`
}

// PerfTimelineMgr provides an interface to maintain a timeline of performance
type PerfTimelineMgr interface {
	Active() bool
	CountStats() int
	Clear()
	Update(stat PerfStat)
	Export() PerfTimeline
}

func copyObj[T any](s *T) *T {
	var ret T
	ret = *s
	return &ret
}

func copyMap[T any](s map[string]T) map[string]T {
	ret := make(map[string]T)
	for k, v := range s {
		ret[k] = v
	}
	return ret
}

func copySlice[T any](s []T) []T {
	ret := make([]T, len(s))
	for i, v := range s {
		ret[i] = v
	}
	return ret
}

func copyPerfStat(s PerfStat) PerfStat {
	return PerfStat{
		CPU:       copyObj(s.CPU),
		Mem:       copyObj(s.Mem),
		NetIOPSec: copyMap(s.NetIOPSec),
		DiskUsage: copyMap(s.DiskUsage),
		CStat:     copyMap(s.CStat),
		CEvent:    copySlice(s.CEvent),
	}
}

type perfTimelineMgr struct {
	ctx      context.Context // context
	interval int64           // interval of each snapshot
	current  PerfStat        // current snapshot
	cinfo    []ContianerInfo // all container info
	datashot []PerfStat      // a list of datashot
	rotate   int             // current index of datashot
	used     int             // used capacity of datashot
}

func CreatePerfTimelineMgr(
	ctx context.Context,
	interval int64,
	capacity int,
	init PerfStat,
) PerfTimelineMgr {
	ret := &perfTimelineMgr{
		ctx:      ctx,
		interval: interval,
		current: PerfStat{
			CPU:       copyObj(init.CPU),
			Mem:       copyObj(init.Mem),
			NetIOPSec: copyMap(init.NetIOPSec),
			DiskUsage: copyMap(init.DiskUsage),
			CStat:     copyMap(init.CStat),
		},
		cinfo: copySlice(init.CEvent),
	}

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ret.ctx.Done():
				return
			case <-ticker.C:
				ret.datashot[ret.rotate] = copyPerfStat(ret.current)
				ret.rotate = (ret.rotate + 1) % capacity
				if ret.used < capacity {
					ret.used++
				}
			}
		}
	}()

	return ret
}
