package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// GetDiskInfo 获取磁盘信息（显示所有挂载点，类似 df -h）
func GetDiskInfo() DiskInfo {
	// true 表示包括所有文件系统，包括 tmpfs、devtmpfs、overlay 等
	partitions, _ := disk.Partitions(true)

	var partitionInfos []PartitionInfo

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		// 过滤掉一些无用的虚拟文件系统（容量为0的）
		// 但保留 tmpfs、overlay 等有实际容量的
		if usage.Total == 0 {
			continue
		}

		partitionInfos = append(partitionInfos, PartitionInfo{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Fstype:      partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	// 按使用率降序排序（使用率高的在前）
	sortDiskByUsage(partitionInfos)

	return DiskInfo{
		Partitions: partitionInfos,
		Timestamp:  time.Now(),
	}
}

// sortDiskByUsage 按磁盘使用率降序排序
func sortDiskByUsage(partitions []PartitionInfo) {
	// 使用冒泡排序（简单实现）
	for i := 0; i < len(partitions)-1; i++ {
		for j := 0; j < len(partitions)-i-1; j++ {
			if partitions[j].UsedPercent < partitions[j+1].UsedPercent {
				partitions[j], partitions[j+1] = partitions[j+1], partitions[j]
			}
		}
	}
}
