package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

// GetNetworkInfo 获取网络信息
func GetNetworkInfo() NetworkInfo {
	ioCounters, _ := net.IOCounters(true)
	interfaces, _ := net.Interfaces()

	var interfaceInfos []InterfaceInfo

	for _, iface := range interfaces {
		// 跳过回环接口
		if iface.Name == "lo" {
			continue
		}

		var addrs []string
		for _, addr := range iface.Addrs {
			addrs = append(addrs, addr.Addr)
		}

		// 查找对应的统计信息
		var bytesSent, bytesRecv, packetsSent, packetsRecv uint64
		for _, counter := range ioCounters {
			if counter.Name == iface.Name {
				bytesSent = counter.BytesSent
				bytesRecv = counter.BytesRecv
				packetsSent = counter.PacketsSent
				packetsRecv = counter.PacketsRecv
				break
			}
		}

		interfaceInfos = append(interfaceInfos, InterfaceInfo{
			Name:        iface.Name,
			BytesSent:   bytesSent,
			BytesRecv:   bytesRecv,
			PacketsSent: packetsSent,
			PacketsRecv: packetsRecv,
			Addrs:       addrs,
		})
	}

	return NetworkInfo{
		Interfaces: interfaceInfos,
		Timestamp:  time.Now(),
	}
}
