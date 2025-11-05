package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// GetDockerInfo 获取 Docker 信息
func GetDockerInfo() DockerInfo {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return DockerInfo{Available: false, Timestamp: time.Now()}
	}
	defer cli.Close()

	// 测试连接
	_, err = cli.Ping(ctx)
	if err != nil {
		return DockerInfo{Available: false, Timestamp: time.Now()}
	}

	// 获取所有容器
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return DockerInfo{Available: false, Timestamp: time.Now()}
	}

	var containerInfos []ContainerInfo
	runningCount := 0

	for _, ctr := range containers {
		info := getContainerInfo(ctx, cli, ctr)
		containerInfos = append(containerInfos, info)

		if info.State == "running" {
			runningCount++
		}
	}

	return DockerInfo{
		Available:    true,
		Containers:   containerInfos,
		RunningCount: runningCount,
		TotalCount:   len(containers),
		Timestamp:    time.Now(),
	}
}

// GetContainerDetail 获取特定容器的详细信息
func GetContainerDetail(containerID string) ContainerInfo {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return ContainerInfo{}
	}
	defer cli.Close()

	// 获取容器列表
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return ContainerInfo{}
	}

	for _, ctr := range containers {
		if ctr.ID == containerID || ctr.ID[:12] == containerID {
			return getContainerInfo(ctx, cli, ctr)
		}
	}

	return ContainerInfo{}
}

func getContainerInfo(ctx context.Context, cli *client.Client, ctr types.Container) ContainerInfo {
	// 获取容器名称（去掉前导 /）
	name := ctr.Names[0]
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	// 计算运行时间
	created := time.Unix(ctr.Created, 0)
	uptime := formatUptime(time.Since(created))

	// 获取端口映射
	var ports []PortMapping
	for _, port := range ctr.Ports {
		mapping := PortMapping{
			PrivatePort: port.PrivatePort,
			PublicPort:  port.PublicPort,
			Type:        port.Type,
			IP:          port.IP,
		}
		ports = append(ports, mapping)
	}

	info := ContainerInfo{
		ID:      ctr.ID[:12],
		Name:    name,
		Image:   ctr.Image,
		Status:  ctr.Status,
		State:   ctr.State,
		Ports:   ports,
		Created: created,
		Uptime:  uptime,
	}

	// 如果容器正在运行，获取统计信息
	if ctr.State == "running" {
		stats, err := cli.ContainerStats(ctx, ctr.ID, false)
		if err == nil {
			defer stats.Body.Close()

			var v types.StatsJSON
			if err := json.NewDecoder(stats.Body).Decode(&v); err == nil {
				// CPU 使用率
				cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage)
				systemDelta := float64(v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage)
				cpuPercent := 0.0
				if systemDelta > 0 && cpuDelta > 0 {
					cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
				}
				info.CPUPercent = cpuPercent

				// 内存使用
				info.MemoryUsageMB = float64(v.MemoryStats.Usage) / 1024 / 1024
				info.MemoryLimitMB = float64(v.MemoryStats.Limit) / 1024 / 1024
				if v.MemoryStats.Limit > 0 {
					info.MemPercent = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100
				}

				// 网络 I/O
				for _, netStats := range v.Networks {
					info.NetInputMB += float64(netStats.RxBytes) / 1024 / 1024
					info.NetOutputMB += float64(netStats.TxBytes) / 1024 / 1024
				}

				// 磁盘 I/O
				for _, bioStats := range v.BlkioStats.IoServiceBytesRecursive {
					if bioStats.Op == "read" {
						info.BlockInputMB += float64(bioStats.Value) / 1024 / 1024
					} else if bioStats.Op == "write" {
						info.BlockOutputMB += float64(bioStats.Value) / 1024 / 1024
					}
				}
			}
		}
	}

	return info
}

func formatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}
