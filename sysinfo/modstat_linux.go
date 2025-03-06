package sysinfo

import "time"

// ProcState represent processor state that read from /proc/[PID]/stat
type ProcState string

// All known PID state list
const (
	ProcSRuning   ProcState = "R" // Running
	ProcSSleeping           = "S" // Sleeping (interruptible wait)
	ProcSWaiting            = "D" // Waiting (uninterruptible wait - usually IO)
	ProcSZombie             = "Z" // Zombie
	ProcSStopped            = "T" // Stopped (debugging)
	ProcSTracing            = "t" // Tracing stop (debugging)
	ProcSPaging             = "W" // Paging (only before Linux 2.6.0)
	ProcSDeadLgc            = "X" // Dead (from Linux 2.6.0 onward)
	ProcSDead               = "x" // Dead (from Linux 2.6.33 to 3.13 only)
	ProcSWakekill           = "K" // Wakekill (from Linux 2.6.33 to 3.13 only)
	ProcSWake               = "W" // Waking (from Linux 2.6.33 to 3.13 only)
	ProcSParking            = "P" // Parked (from Linux 3.9 to 3.13 only)
	ProcSIdle               = "I" // Idle (from Linux 4.14 onward)
	ProcSUnknow             = "-" // Unknow
)

// A PIDInfo is a short report for a process that read from /proc/[PID]/stat
type PIDInfo struct {
	PID    int       `json:"pid"`
	Comma  string    `json:"comma"`
	Cmd    []string  `json:"cmd"`
	Status ProcState `json:"status"`
	PPID   int       `json:"ppid"`
	PGroup int       `json:"pgroup"`
	Detail *PIDStat  `json:"detail,omitempty"`
	NetIO  *PIDNetIO `json:"net_io,omitempty"`
	BlkIO  *PIDBlkIO `json:"block_io,omitempty"`
}

// PIDNetIO is netwotk stream short report for a net work device.
// it read from /proc/net/dev
type PIDNetIO struct {
	ByteRead  uint64 // read bytes from device
	ByteWrite uint64 // write bytes to device
	PackRead  uint64 // read packages from device
	PackWrite uint64 // write packages to device
}

// A PIDNetStat contian process network staticstics data
// It's a mapping of the SNMP. also, will get IO statistics too.
// reference: man/proc(5)
type PIDNetStat struct {
	// statistics IO stream for each device
	NetIO     map[string]PIDNetIO `json:"if_stat"`
	NumConn   uint                `json:"conn"`      // count connections
	InDGrams  uint64              `json:"in_dgram"`  // received data grams
	OutDGrams uint64              `json:"out_dgram"` // sent data grams
}

// A PIDBlkIO contian block IO statistics
// reference: man/proc(5)
// reference: man/ptrace(2)
type PIDBlkIO struct {
	// count charaters readed from character device
	CharRead uint64 `json:"r_char"`
	// count charaters writed to character device
	CharWrite uint64 `json:"w_char"`
	// count systemcall bytes readed
	SyscRead uint64 `json:"r_sysc"`
	// count systemcall bytes writed
	SyscWrite uint64 `json:"w_sysc"`
	// count bytes readed from non-character device
	ByteRead uint64 `json:"r_byte"`
	// count bytes writed to non-character device
	ByteWrite uint64 `json:"w_byte"`
}

// A PIDStat contian PID detail from /proc/[PID]/stat
// it was parsed to golang data type
// reference: man/proc(5)
type PIDStat struct {
	TTYnr      int           `json:"tty_nr"`       // 7
	TpGID      int           `json:"tpgid"`        // 8
	Flags      uint          `json:"falgs"`        // 9
	MinFlt     uint64        `json:"min_flt"`      // 10
	CMinFlt    uint64        `json:"chd_min_flt"`  // 11
	MajFlt     uint64        `json:"maj_flt"`      // 12
	CMajFlt    uint64        `json:"chd_maj_flt"`  // 13
	UTime      time.Duration `json:"usr_time"`     // 14
	STime      time.Duration `json:"sys_time"`     // 15
	CUTime     time.Duration `json:"chd_usr_time"` // 16
	CSTime     time.Duration `json:"chd_sys_time"` // 17
	Priority   int64         `json:"priority"`     // 18
	Nice       int64         `json:"nice"`         // 19
	NumThreads int64         `json:"num_threads"`  // 20
	// Obsolete 21
	StartTime  time.Duration `json:"start_time"`  // 22
	VSize      uint64        `json:"vm_size"`     // 23 num bytes
	RSS        int64         `json:"rss"`         // 24 num bytes
	RSSLim     uint64        `json:"limit_rss"`   // 25
	StartCode  uint          `json:"addr_start"`  // 26
	EndCode    uint          `json:"addr_end"`    // 27
	StartStack uint          `json:"addr_stack"`  // 28
	KStkESP    uint          `json:"current_esp"` // 29
	KStkEIP    uint          `json:"current_eip"` // 30
	// Obsolete 31 - 34
	WChan      uint          `json:"addr_wchan"`     // 35
	NSwap      uint          `json:"num_swaped"`     // 36
	CNSwap     uint          `json:"chd_num_swaped"` // 37
	ExitSignal int           `json:"exit_sig"`       // 38
	Processor  int           `json:"processor"`      // 39
	RTPriority uint          `json:"rt_priority"`    // 40
	Policy     uint          `json:"policy"`         // 41
	DelayBlkIO time.Duration `json:"blkio_delay"`    // 42 delayacct_blkio_ticks
	GuestTime  time.Duration `json:"guest_time"`     // 43
	CGuestTime time.Duration `json:"chd_guest_time"` // 44
	StartData  uint64        `json:"addr_dat_start"` // 45
	EndData    uint64        `json:"addr_dat_end"`   // 46
	StartBrk   uint64        `json:"addr_brk"`       // 47
	ArgStart   uint64        `json:"arg_start"`      // 48
	ArgEnd     uint64        `json:"arg_end"`        // 49
	EnvStart   uint64        `json:"env_start"`      // 50
	EnvEnd     uint64        `json:"env_end"`        // 51
	ExitCode   int           `json:"exitcode"`       // 52
}

// String return the string representation of the PIDStat
func (s ProcState) String() string {
	switch s {
	case ProcSRuning:
		return "Running"
	case ProcSSleeping:
		return "Sleeping"
	case ProcSWaiting:
		return "Waiting"
	case ProcSZombie:
		return "Zombie"
	case ProcSStopped:
		return "Stopped"
	case ProcSTracing:
		return "Tracing"
	case ProcSPaging:
		return "Paging"
	case ProcSDeadLgc:
		return "Dead"
	case ProcSDead:
		return "Dead"
	case ProcSWakekill:
		return "Wakekill"
	case ProcSWake:
		return "Waking"
	case ProcSParking:
		return "Parked"
	case ProcSIdle:
		return "Idle"
	default:
		return "Unknow"
	}
}
