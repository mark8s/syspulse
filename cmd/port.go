package cmd

import (
	"fmt"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "显示端口监听信息",
	Long:  "显示系统正在监听的端口和对应的进程",
	Run: func(cmd *cobra.Command, args []string) {
		display.Clear()
		display.PrintHeader("🔌 端口监听信息")

		portInfo := monitor.GetPortInfo()
		display.PrintPortInfo(portInfo)

		fmt.Println()
		display.PrintFooter("数据更新时间: " + portInfo.Timestamp.Format("2006-01-02 15:04:05"))
	},
}
