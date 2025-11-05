package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "syspulse",
	Short: "🚀 SysPulse - 超级易用的 Linux 系统资源监控工具",
	Long: `SysPulse 是一个美观、直观的系统资源监控工具
	
支持监控:
  • CPU、内存、磁盘、网络
  • Docker 容器资源占用
  • 进程信息
  • 实时刷新模式`,
	Run: func(cmd *cobra.Command, args []string) {
		// 默认显示仪表盘
		dashboardCmd.Run(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
	rootCmd.AddCommand(cpuCmd)
	rootCmd.AddCommand(memoryCmd)
	rootCmd.AddCommand(diskCmd)
	rootCmd.AddCommand(networkCmd)
	rootCmd.AddCommand(portCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(dockerCmd)
	rootCmd.AddCommand(webCmd)
}
