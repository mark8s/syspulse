package cmd

import (
	"fmt"
	"time"

	"syspulse/internal/display"
	"syspulse/internal/monitor"

	"github.com/spf13/cobra"
)

var (
	dockerWatch    bool
	dockerInterval int
	containerID    string
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "显示 Docker 容器信息",
	Long:  "显示 Docker 容器运行状态和资源占用",
	Run: func(cmd *cobra.Command, args []string) {
		if dockerWatch {
			runDockerWatchMode()
		} else {
			showDockerInfo()
		}
	},
}

func init() {
	dockerCmd.Flags().BoolVarP(&dockerWatch, "watch", "w", false, "实时刷新模式")
	dockerCmd.Flags().IntVarP(&dockerInterval, "interval", "i", 2, "刷新间隔（秒）")
	dockerCmd.Flags().StringVarP(&containerID, "container", "c", "", "查看特定容器详情")
}

func showDockerInfo() {
	display.Clear()
	display.PrintHeader("🐳 Docker 容器监控")

	dockerInfo := monitor.GetDockerInfo()

	if !dockerInfo.Available {
		display.PrintError("❌ Docker 不可用")
		fmt.Println("   请确保:")
		fmt.Println("   1. Docker 已安装")
		fmt.Println("   2. Docker 服务正在运行")
		fmt.Println("   3. 当前用户有 Docker 权限")
		return
	}

	if containerID != "" {
		// 显示特定容器的详细信息
		containerInfo := monitor.GetContainerDetail(containerID)
		display.PrintContainerDetail(containerInfo)
	} else {
		// 显示所有容器概览
		display.PrintDockerInfoDetailed(dockerInfo)
	}

	fmt.Println()
	display.PrintFooter("数据更新时间: " + dockerInfo.Timestamp.Format("2006-01-02 15:04:05"))
}

func runDockerWatchMode() {
	ticker := time.NewTicker(time.Duration(dockerInterval) * time.Second)
	defer ticker.Stop()

	// 立即显示第一次
	showDockerInfo()

	for range ticker.C {
		showDockerInfo()
	}
}
