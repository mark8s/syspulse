package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

// GetMemoryInfo 获取内存信息
func GetMemoryInfo() MemoryInfo {
	vmem, _ := mem.VirtualMemory()
	swap, _ := mem.SwapMemory()

	return MemoryInfo{
		Total:       vmem.Total,
		Used:        vmem.Used,
		Available:   vmem.Available,
		UsedPercent: vmem.UsedPercent,
		SwapTotal:   swap.Total,
		SwapUsed:    swap.Used,
		SwapPercent: swap.UsedPercent,
		Cached:      vmem.Cached,
		Buffers:     vmem.Buffers,
		Timestamp:   time.Now(),
	}
}
