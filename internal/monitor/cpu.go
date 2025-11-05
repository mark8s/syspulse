package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

// GetCPUInfo 获取 CPU 信息
func GetCPUInfo() CPUInfo {
	// CPU 使用率
	percent, _ := cpu.Percent(time.Second, false)
	cpuPercent := 0.0
	if len(percent) > 0 {
		cpuPercent = percent[0]
	}

	// 每个核心的使用率
	perCorePercent, _ := cpu.Percent(time.Second, true)

	// CPU 核心数
	coreCount, _ := cpu.Counts(true)

	// CPU 信息
	cpuInfos, _ := cpu.Info()
	modelName := "Unknown"
	if len(cpuInfos) > 0 {
		modelName = cpuInfos[0].ModelName
	}

	// 负载平均值
	loadAvg, _ := load.Avg()

	return CPUInfo{
		UsagePercent: cpuPercent,
		CoreCount:    coreCount,
		ModelName:    modelName,
		LoadAvg1:     loadAvg.Load1,
		LoadAvg5:     loadAvg.Load5,
		LoadAvg15:    loadAvg.Load15,
		PerCoreUsage: perCorePercent,
		Timestamp:    time.Now(),
	}
}
