package cmd

import (
	"fmt"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "显示 CPU 信息",
	Long:  "显示详细的 CPU 使用率、核心数、负载等信息",
	Run: func(cmd *cobra.Command, args []string) {
		display.Clear()
		display.PrintHeader("🔥 CPU 信息")

		cpuInfo := monitor.GetCPUInfo()
		display.PrintCPUInfoDetailed(cpuInfo)

		fmt.Println()
		display.PrintFooter("数据更新时间: " + cpuInfo.Timestamp.Format("2006-01-02 15:04:05"))
	},
}
