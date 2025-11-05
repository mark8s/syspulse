package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// PortInfo 端口信息
type PortInfo struct {
	Listening []PortDetail
	Timestamp time.Time
}

// PortDetail 端口详情
type PortDetail struct {
	Port        uint32
	Address     string
	Protocol    string
	State       string
	PID         int32
	ProcessName string
}

// GetPortInfo 获取端口信息
func GetPortInfo() PortInfo {
	connections, _ := net.Connections("all")

	portMap := make(map[string]*PortDetail)

	for _, conn := range connections {
		// 只关注 LISTEN 状态的连接
		if conn.Status != "LISTEN" {
			continue
		}

		key := fmt.Sprintf("%s:%d", conn.Type, conn.Laddr.Port)

		// 避免重复
		if _, exists := portMap[key]; exists {
			continue
		}

		detail := &PortDetail{
			Port:     conn.Laddr.Port,
			Address:  conn.Laddr.IP,
			Protocol: getProtocolName(conn.Type),
			State:    conn.Status,
			PID:      conn.Pid,
		}

		// 尝试获取进程名
		if conn.Pid > 0 {
			if p, err := process.NewProcess(conn.Pid); err == nil {
				if name, err := p.Name(); err == nil {
					detail.ProcessName = name
				}
			}
		}

		portMap[key] = detail
	}

	// 转换为切片
	var listening []PortDetail
	for _, detail := range portMap {
		listening = append(listening, *detail)
	}

	return PortInfo{
		Listening: listening,
		Timestamp: time.Now(),
	}
}

func getProtocolName(connType uint32) string {
	switch connType {
	case 1:
		return "tcp"
	case 2:
		return "udp"
	case 3:
		return "tcp6"
	case 4:
		return "udp6"
	default:
		return "unknown"
	}
}
