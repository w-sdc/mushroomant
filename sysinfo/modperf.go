package sysinfo

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrPerfTimelineMgrClosed = errors.New("perf timeline manager is closed")
)

// CPUStat represents the CPU usage of the system
type CPUStat struct {
	Total float32   `json:"total"`
	Core  []float32 `json:"core"`
}

// MemStat represents the memory usage of the system
type MemStat struct {
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
	Available uint64 `json:"available"`
}

// NetStat represents the network I/O of the system
type NetStat struct {
	BytesSend   uint64 `json:"bytes_send"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSend uint64 `json:"packets_send"`
	PacketsRecv uint64 `json:"packets_recv"`
}

// DiskUsage represents the disk usage of the system
type DiskUsage struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

// ContainerStat represents the performance data of a container
type ContainerStat struct {
	CPU     float32 `json:"cpu"`
	MemUsed uint64  `json:"mem_used"`
}

// ContianerInfo represents the information of a container
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

// PerfTimeline is a struct that contains a timeline of performance data
type PerfTimeline struct {
	Interval int64 `json:"interval"`
	// Stats is a list of performance data, in descending order of time
	Stats []PerfStat               `json:"stats"`
	CInfo map[string]ContianerInfo `json:"c_event"`
}

// PerfTimelineMgr provides an interface to maintain a timeline of performance
type PerfTimelineMgr interface {
	Active() bool
	CountStats() int
	Clear()
	Update(stat PerfStat) error
	Export() PerfTimeline
}

// copyObj is a generic function to copy an pointer object
func copyObj[T any](s *T) *T {
	var ret T
	ret = *s
	return &ret
}

// copyMap is a generic function to copy a map
func copyMap[T any](s map[string]T) map[string]T {
	ret := make(map[string]T)
	for k, v := range s {
		ret[k] = v
	}
	return ret
}

// copySlice is a generic function to copy a slice
func copySlice[T any](s []T) []T {
	ret := make([]T, len(s))
	for i, v := range s {
		ret[i] = v
	}
	return ret
}

// copyPerfStat is used to copy a PerfStat object
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

// perfTimelineMgr is an implementation of PerfTimelineMgr
type perfTimelineMgr struct {
	ctx      context.Context          // context
	mtx      sync.RWMutex             // mutex
	interval int64                    // interval of each snapshot
	current  PerfStat                 // current snapshot
	cinfo    map[string]ContianerInfo // all container info
	datashot []PerfStat               // a list of datashot
	rotate   int                      // current index of datashot
	used     int                      // used capacity of datashot
}

// CreatePerfTimelineMgr creates a PerfTimelineMgr object
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
		cinfo: make(map[string]ContianerInfo),
	}
	for _, v := range init.CEvent {
		ret.cinfo[v.ID] = v
	}

	// start time line snapshot
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ret.ctx.Done():
				return
			case <-ticker.C:
				func() {
					ret.mtx.Lock()
					defer ret.mtx.Unlock()
					ret.datashot[ret.rotate] = copyPerfStat(ret.current)
					ret.rotate = (ret.rotate + 1) % capacity
					if ret.used < capacity {
						ret.used++
					}
				}()
			}
		}
	}()

	return ret
}

// Active returns true if the manager is still active
func (m *perfTimelineMgr) Active() bool {
	select {
	case <-m.ctx.Done():
		return false
	default:
		return true
	}
}

// CountStats returns the number of stats in the timeline
func (m *perfTimelineMgr) CountStats() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.used
}

// Clear clears all the data in the timeline
func (m *perfTimelineMgr) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.used = 0
	m.rotate = 0
}

// Update updates the current snapshot
func (m *perfTimelineMgr) Update(stat PerfStat) error {
	if !m.Active() {
		return ErrPerfTimelineMgrClosed
	}
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if stat.CPU != nil {
		m.current.CPU = copyObj(stat.CPU)
	}
	if stat.Mem != nil {
		m.current.Mem = copyObj(stat.Mem)
	}
	if stat.NetIOPSec != nil {
		m.current.NetIOPSec = copyMap(stat.NetIOPSec)
	}
	if stat.DiskUsage != nil {
		m.current.DiskUsage = copyMap(stat.DiskUsage)
	}
	if stat.CStat != nil {
		m.current.CStat = copyMap(stat.CStat)
	}
	if stat.CEvent != nil {
		m.current.CEvent = copySlice(stat.CEvent)
		for _, v := range stat.CEvent {
			m.cinfo[v.ID] = v
		}
	}
	return nil
}

// Export exports the timeline data
func (m *perfTimelineMgr) Export() PerfTimeline {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	ret := PerfTimeline{
		Interval: m.interval,
		Stats:    make([]PerfStat, m.used),
		CInfo:    copyMap(m.cinfo),
	}
	for i := 0; i < m.used; i++ {
		ret.Stats[i] = copyPerfStat(m.datashot[(m.rotate+i)%len(m.datashot)])
	}
	return ret
}
