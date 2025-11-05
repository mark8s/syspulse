package monitor

import (
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// GetProcessInfo 获取进程信息
func GetProcessInfo(topN int) ProcessInfo {
	processes, _ := process.Processes()

	var processDetails []ProcessDetail

	for _, p := range processes {
		name, _ := p.Name()
		username, _ := p.Username()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()
		status, _ := p.Status()
		cmdline, _ := p.Cmdline()

		memoryMB := float64(0)
		if memInfo != nil {
			memoryMB = float64(memInfo.RSS) / 1024 / 1024
		}

		processDetails = append(processDetails, ProcessDetail{
			PID:        p.Pid,
			Name:       name,
			Username:   username,
			CPUPercent: cpuPercent,
			MemoryMB:   memoryMB,
			MemPercent: memPercent,
			Status:     status[0],
			Command:    cmdline,
		})
	}

	// 按 CPU 排序
	topCPU := make([]ProcessDetail, len(processDetails))
	copy(topCPU, processDetails)
	sort.Slice(topCPU, func(i, j int) bool {
		return topCPU[i].CPUPercent > topCPU[j].CPUPercent
	})
	if len(topCPU) > topN {
		topCPU = topCPU[:topN]
	}

	// 按内存排序
	topMemory := make([]ProcessDetail, len(processDetails))
	copy(topMemory, processDetails)
	sort.Slice(topMemory, func(i, j int) bool {
		return topMemory[i].MemoryMB > topMemory[j].MemoryMB
	})
	if len(topMemory) > topN {
		topMemory = topMemory[:topN]
	}

	return ProcessInfo{
		TotalProcesses: len(processes),
		TopCPU:         topCPU,
		TopMemory:      topMemory,
		Timestamp:      time.Now(),
	}
}
