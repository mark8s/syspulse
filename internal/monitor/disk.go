package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// GetDiskInfo 获取磁盘信息
func GetDiskInfo() DiskInfo {
	partitions, _ := disk.Partitions(false)

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
