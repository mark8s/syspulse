package cmd

import (
	"fmt"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var (
	topN int
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "显示进程信息",
	Long:  "显示系统进程列表和资源占用 Top N",
	Run: func(cmd *cobra.Command, args []string) {
		display.Clear()
		display.PrintHeader("⚙️  进程信息")

		processInfo := monitor.GetProcessInfo(topN)
		display.PrintProcessInfo(processInfo)

		fmt.Println()
		display.PrintFooter("数据更新时间: " + processInfo.Timestamp.Format("2006-01-02 15:04:05"))
	},
}

func init() {
	processCmd.Flags().IntVarP(&topN, "top", "t", 10, "显示 Top N 进程")
}
