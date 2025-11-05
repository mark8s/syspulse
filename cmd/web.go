package cmd

import (
	"fmt"

	"syspulse/internal/web"

	"github.com/spf13/cobra"
)

var (
	webPort int
	webHost string
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动 Web 界面服务器",
	Long:  "启动一个 Web 服务器，通过浏览器查看系统资源监控",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🌐 正在启动 SysPulse Web 服务器...\n")
		fmt.Printf("📡 地址: http://%s:%d\n", webHost, webPort)
		fmt.Printf("💡 在浏览器中打开上面的地址即可查看监控面板\n")
		fmt.Printf("⏹️  按 Ctrl+C 停止服务器\n\n")

		server := web.NewServer(webHost, webPort)
		if err := server.Start(); err != nil {
			fmt.Printf("❌ 启动失败: %v\n", err)
		}
	},
}

func init() {
	webCmd.Flags().IntVarP(&webPort, "port", "p", 3000, "Web 服务器端口")
	webCmd.Flags().StringVarP(&webHost, "host", "H", "0.0.0.0", "Web 服务器主机地址")
}
