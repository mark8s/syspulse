package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

// GetSystemInfo 获取系统基本信息
func GetSystemInfo() SystemInfo {
	info, _ := host.Info()

	return SystemInfo{
		Hostname:  info.Hostname,
		OS:        info.OS,
		Kernel:    info.KernelVersion,
		Uptime:    info.Uptime,
		Timestamp: time.Now(),
	}
}
