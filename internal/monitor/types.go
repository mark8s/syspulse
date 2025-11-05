package monitor

import "time"

// SystemInfo 系统基本信息
type SystemInfo struct {
	Hostname  string
	OS        string
	Kernel    string
	Uptime    uint64
	Timestamp time.Time
}

// CPUInfo CPU 信息
type CPUInfo struct {
	UsagePercent float64
	CoreCount    int
	ModelName    string
	LoadAvg1     float64
	LoadAvg5     float64
	LoadAvg15    float64
	PerCoreUsage []float64
	Timestamp    time.Time
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total       uint64
	Used        uint64
	Available   uint64
	UsedPercent float64
	SwapTotal   uint64
	SwapUsed    uint64
	SwapPercent float64
	Cached      uint64
	Buffers     uint64
	Timestamp   time.Time
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Partitions []PartitionInfo
	Timestamp  time.Time
}

// PartitionInfo 分区信息
type PartitionInfo struct {
	Device      string
	Mountpoint  string
	Fstype      string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Interfaces []InterfaceInfo
	Timestamp  time.Time
}

// InterfaceInfo 网络接口信息
type InterfaceInfo struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	Addrs       []string
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	TotalProcesses int
	TopCPU         []ProcessDetail
	TopMemory      []ProcessDetail
	Timestamp      time.Time
}

// ProcessDetail 进程详情
type ProcessDetail struct {
	PID        int32
	Name       string
	Username   string
	CPUPercent float64
	MemoryMB   float64
	MemPercent float32
	Status     string
	Command    string
}

// DockerInfo Docker 信息
type DockerInfo struct {
	Available    bool
	Containers   []ContainerInfo
	RunningCount int
	TotalCount   int
	Timestamp    time.Time
}

// ContainerInfo 容器信息
type ContainerInfo struct {
	ID            string
	Name          string
	Image         string
	Status        string
	State         string
	Ports         []PortMapping
	CPUPercent    float64
	MemoryUsageMB float64
	MemoryLimitMB float64
	MemPercent    float64
	NetInputMB    float64
	NetOutputMB   float64
	BlockInputMB  float64
	BlockOutputMB float64
	Created       time.Time
	Uptime        string
}

// PortMapping 端口映射
type PortMapping struct {
	PrivatePort uint16
	PublicPort  uint16
	Type        string
	IP          string
}
