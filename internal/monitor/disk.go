package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// GetDiskInfo 获取磁盘信息（显示所有挂载点，包括虚拟文件系统）
func GetDiskInfo() DiskInfo {
	// true 表示包括所有文件系统，包括 tmpfs、devtmpfs、overlay 等
	partitions, _ := disk.Partitions(true)

	var partitionInfos []PartitionInfo

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
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

	return DiskInfo{
		Partitions: partitionInfos,
		Timestamp:  time.Now(),
	}
}
