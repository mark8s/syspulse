package cmd

import (
	"fmt"
	"os"
	"time"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var (
	watchMode bool
	interval  int
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "显示系统资源仪表盘",
	Long:  "显示包含 CPU、内存、磁盘、网络、Docker 容器的完整仪表盘",
	Run: func(cmd *cobra.Command, args []string) {
		if watchMode {
			runWatchMode()
		} else {
			showDashboard()
		}
	},
}

func init() {
	dashboardCmd.Flags().BoolVarP(&watchMode, "watch", "w", false, "实时刷新模式")
	dashboardCmd.Flags().IntVarP(&interval, "interval", "i", 2, "刷新间隔（秒）")
}

func showDashboard() {
	// 清屏
	display.Clear()

	// 显示标题
	display.PrintHeader("💻 SYSTEM PULSE - 系统概览")

	// 获取系统信息
	sysInfo := monitor.GetSystemInfo()
	display.PrintSystemInfo(sysInfo)

	fmt.Println()

	// CPU 信息
	cpuInfo := monitor.GetCPUInfo()
	display.PrintCPUInfo(cpuInfo)

	fmt.Println()

	// 内存信息
	memInfo := monitor.GetMemoryInfo()
	display.PrintMemoryInfo(memInfo)

	fmt.Println()

	// 磁盘信息
	diskInfo := monitor.GetDiskInfo()
	display.PrintDiskInfo(diskInfo)

	fmt.Println()

	// 网络信息
	netInfo := monitor.GetNetworkInfo()
	display.PrintNetworkInfo(netInfo)

	fmt.Println()

	// Docker 容器信息
	dockerInfo := monitor.GetDockerInfo()
	if dockerInfo.Available {
		display.PrintDockerInfo(dockerInfo)
	} else {
		display.PrintWarning("🐳 Docker 不可用或未运行")
	}

	fmt.Println()
	display.PrintFooter("按 Ctrl+C 退出")
}

func runWatchMode() {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// 立即显示第一次
	showDashboard()

	for range ticker.C {
		showDashboard()
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	// Windows 兼容
	if os.Getenv("OS") == "Windows_NT" {
		fmt.Print("\033c")
	}
}
