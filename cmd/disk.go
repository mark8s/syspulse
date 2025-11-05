package cmd

import (
	"fmt"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "显示磁盘信息 (类似 df -h)",
	Long:  "显示所有磁盘分区的使用情况，输出格式类似 df -h 命令",
	Run: func(cmd *cobra.Command, args []string) {
		display.Clear()
		display.PrintHeader("💿 磁盘使用情况 (df -h)")

		diskInfo := monitor.GetDiskInfo()
		display.PrintDiskInfoDetailed(diskInfo)

		fmt.Println()
		display.PrintFooter("数据更新时间: " + diskInfo.Timestamp.Format("2006-01-02 15:04:05"))
	},
}
