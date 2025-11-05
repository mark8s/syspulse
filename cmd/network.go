package cmd

import (
	"fmt"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "显示网络信息",
	Long:  "显示网络接口、流量统计和连接信息",
	Run: func(cmd *cobra.Command, args []string) {
		display.Clear()
		display.PrintHeader("🌐 网络信息")

		netInfo := monitor.GetNetworkInfo()
		display.PrintNetworkInfoDetailed(netInfo)

		fmt.Println()
		display.PrintFooter("数据更新时间: " + netInfo.Timestamp.Format("2006-01-02 15:04:05"))
	},
}
