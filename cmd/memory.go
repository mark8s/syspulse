package cmd

import (
	"fmt"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "显示内存信息",
	Long:  "显示详细的内存使用情况，包括物理内存和 Swap",
	Run: func(cmd *cobra.Command, args []string) {
		display.Clear()
		display.PrintHeader("💾 内存信息")

		memInfo := monitor.GetMemoryInfo()
		display.PrintMemoryInfoDetailed(memInfo)

		fmt.Println()
		display.PrintFooter("数据更新时间: " + memInfo.Timestamp.Format("2006-01-02 15:04:05"))
	},
}
