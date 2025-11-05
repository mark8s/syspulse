package display

import (
	"fmt"
	"os"
	"strings"

	"syspulse/internal/monitor"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	// 颜色定义
	colorTitle   = color.New(color.FgCyan, color.Bold)
	colorSuccess = color.New(color.FgGreen)
	colorWarning = color.New(color.FgYellow)
	colorError   = color.New(color.FgRed)
	colorInfo    = color.New(color.FgBlue)
	colorValue   = color.New(color.FgWhite, color.Bold)
	colorLabel   = color.New(color.FgHiBlack)
)

// Clear 清屏
func Clear() {
	fmt.Print("\033[H\033[2J")
}

// PrintHeader 打印标题
func PrintHeader(title string) {
	width := 80
	border := strings.Repeat("═", width-4)

	fmt.Println()
	colorTitle.Println("╔" + border + "╗")
	padding := (width - len(title) - 4) / 2
	colorTitle.Printf("║%s%s%s║\n",
		strings.Repeat(" ", padding),
		title,
		strings.Repeat(" ", width-4-padding-len(title)))
	colorTitle.Println("╚" + border + "╝")
	fmt.Println()
}

// PrintFooter 打印页脚
func PrintFooter(text string) {
	colorLabel.Println(strings.Repeat("─", 80))
	colorInfo.Println(text)
}

// PrintWarning 打印警告
func PrintWarning(text string) {
	colorWarning.Println(text)
}

// PrintError 打印错误
func PrintError(text string) {
	colorError.Println(text)
}

// PrintSystemInfo 打印系统信息
func PrintSystemInfo(info monitor.SystemInfo) {
	uptime := formatUptime(info.Uptime)

	colorTitle.Println("🖥️  系统信息")
	fmt.Printf("  ")
	colorLabel.Print("主机名: ")
	colorValue.Println(info.Hostname)

	fmt.Printf("  ")
	colorLabel.Print("操作系统: ")
	colorValue.Println(info.OS)

	fmt.Printf("  ")
	colorLabel.Print("内核版本: ")
	colorValue.Println(info.Kernel)

	fmt.Printf("  ")
	colorLabel.Print("运行时长: ")
	colorSuccess.Println(uptime)
}

// PrintCPUInfo 打印 CPU 信息（简洁版）
func PrintCPUInfo(info monitor.CPUInfo) {
	colorTitle.Println("🔥 CPU")

	// 使用率和进度条
	fmt.Printf("  ")
	colorLabel.Print("使用率: ")
	printPercentWithBar(info.UsagePercent, 40)

	fmt.Printf("  ")
	colorLabel.Print("核心数: ")
	colorValue.Println(info.CoreCount)

	fmt.Printf("  ")
	colorLabel.Print("负载: ")
	colorValue.Printf("%.2f / %.2f / %.2f\n", info.LoadAvg1, info.LoadAvg5, info.LoadAvg15)
}

// PrintCPUInfoDetailed 打印 CPU 详细信息
func PrintCPUInfoDetailed(info monitor.CPUInfo) {
	fmt.Printf("  ")
	colorLabel.Print("型号: ")
	colorValue.Println(info.ModelName)

	fmt.Println()
	fmt.Printf("  ")
	colorLabel.Print("总体使用率: ")
	printPercentWithBar(info.UsagePercent, 50)

	fmt.Println()
	fmt.Printf("  ")
	colorLabel.Println("各核心使用率:")

	for i, usage := range info.PerCoreUsage {
		fmt.Printf("    ")
		colorLabel.Printf("核心 %2d: ", i)
		printPercentWithBar(usage, 40)
	}

	fmt.Println()
	fmt.Printf("  ")
	colorLabel.Println("负载平均值:")
	fmt.Printf("    ")
	colorLabel.Print("1 分钟:  ")
	printLoadValue(info.LoadAvg1, info.CoreCount)
	fmt.Printf("    ")
	colorLabel.Print("5 分钟:  ")
	printLoadValue(info.LoadAvg5, info.CoreCount)
	fmt.Printf("    ")
	colorLabel.Print("15 分钟: ")
	printLoadValue(info.LoadAvg15, info.CoreCount)
}

// PrintMemoryInfo 打印内存信息（简洁版）
func PrintMemoryInfo(info monitor.MemoryInfo) {
	colorTitle.Println("💾 内存")

	fmt.Printf("  ")
	colorLabel.Print("物理内存: ")
	colorValue.Printf("%s / %s ", formatBytes(info.Used), formatBytes(info.Total))
	printPercentWithBar(info.UsedPercent, 30)

	fmt.Printf("  ")
	colorLabel.Print("可用内存: ")
	colorSuccess.Println(formatBytes(info.Available))

	if info.SwapTotal > 0 {
		fmt.Printf("  ")
		colorLabel.Print("Swap: ")
		colorValue.Printf("%s / %s ", formatBytes(info.SwapUsed), formatBytes(info.SwapTotal))
		printPercentWithBar(info.SwapPercent, 30)
	}
}

// PrintMemoryInfoDetailed 打印内存详细信息
func PrintMemoryInfoDetailed(info monitor.MemoryInfo) {
	fmt.Printf("  ")
	colorLabel.Println("物理内存:")
	fmt.Printf("    ")
	colorLabel.Print("总量: ")
	colorValue.Println(formatBytes(info.Total))
	fmt.Printf("    ")
	colorLabel.Print("已用: ")
	colorValue.Println(formatBytes(info.Used))
	fmt.Printf("    ")
	colorLabel.Print("可用: ")
	colorSuccess.Println(formatBytes(info.Available))
	fmt.Printf("    ")
	colorLabel.Print("缓存: ")
	colorInfo.Println(formatBytes(info.Cached))
	fmt.Printf("    ")
	colorLabel.Print("缓冲: ")
	colorInfo.Println(formatBytes(info.Buffers))
	fmt.Printf("    ")
	colorLabel.Print("使用率: ")
	printPercentWithBar(info.UsedPercent, 50)

	fmt.Println()
	if info.SwapTotal > 0 {
		fmt.Printf("  ")
		colorLabel.Println("交换分区 (Swap):")
		fmt.Printf("    ")
		colorLabel.Print("总量: ")
		colorValue.Println(formatBytes(info.SwapTotal))
		fmt.Printf("    ")
		colorLabel.Print("已用: ")
		colorValue.Println(formatBytes(info.SwapUsed))
		fmt.Printf("    ")
		colorLabel.Print("使用率: ")
		printPercentWithBar(info.SwapPercent, 50)
	} else {
		fmt.Printf("  ")
		colorLabel.Println("交换分区 (Swap): 未配置")
	}
}

// PrintDiskInfo 打印磁盘信息（简洁版，类似 df -h）
func PrintDiskInfo(info monitor.DiskInfo) {
	colorTitle.Println("💿 磁盘 (df -h)")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"文件系统", "容量", "已用", "可用", "已用%", "挂载点"})
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
	})

	for _, partition := range info.Partitions {
		table.Append([]string{
			partition.Device,
			formatBytes(partition.Total),
			formatBytes(partition.Used),
			formatBytes(partition.Free),
			fmt.Sprintf("%.0f%%", partition.UsedPercent),
			partition.Mountpoint,
		})
	}

	table.Render()
}

// PrintDiskInfoDetailed 打印磁盘详细信息（完整版 df -h）
func PrintDiskInfoDetailed(info monitor.DiskInfo) {
	fmt.Println()
	colorTitle.Println("文件系统磁盘使用情况 (Filesystem disk space usage)")
	fmt.Println()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"文件系统", "类型", "容量", "已用", "可用", "已用%", "挂载点"})
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
	})

	for _, partition := range info.Partitions {
		// 根据使用率设置颜色提示
		usageStr := fmt.Sprintf("%.0f%%", partition.UsedPercent)

		table.Append([]string{
			partition.Device,
			partition.Fstype,
			formatBytes(partition.Total),
			formatBytes(partition.Used),
			formatBytes(partition.Free),
			usageStr,
			partition.Mountpoint,
		})
	}

	table.Render()

	// 添加警告提示
	fmt.Println()
	hasWarning := false
	for _, partition := range info.Partitions {
		if partition.UsedPercent >= 90 {
			colorError.Printf("⚠️  警告: %s 使用率已达 %.0f%%，空间不足！\n", partition.Mountpoint, partition.UsedPercent)
			hasWarning = true
		} else if partition.UsedPercent >= 80 {
			colorWarning.Printf("⚠️  提示: %s 使用率已达 %.0f%%，建议清理\n", partition.Mountpoint, partition.UsedPercent)
			hasWarning = true
		}
	}

	if !hasWarning {
		colorSuccess.Println("✅ 所有磁盘空间充足")
	}
}

// PrintNetworkInfo 打印网络信息（简洁版）
func PrintNetworkInfo(info monitor.NetworkInfo) {
	colorTitle.Println("🌐 网络")

	for _, iface := range info.Interfaces {
		if len(iface.Addrs) == 0 {
			continue
		}

		fmt.Printf("  ")
		colorLabel.Printf("%s: ", iface.Name)
		colorValue.Println(iface.Addrs[0])

		fmt.Printf("    ")
		colorLabel.Print("↑ ")
		colorSuccess.Printf("%s  ", formatBytes(iface.BytesSent))
		colorLabel.Print("↓ ")
		colorInfo.Println(formatBytes(iface.BytesRecv))
	}
}

// PrintNetworkInfoDetailed 打印网络详细信息
func PrintNetworkInfoDetailed(info monitor.NetworkInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"接口", "地址", "发送", "接收", "发送包", "接收包"})
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, iface := range info.Interfaces {
		addr := "-"
		if len(iface.Addrs) > 0 {
			addr = iface.Addrs[0]
		}

		table.Append([]string{
			iface.Name,
			addr,
			formatBytes(iface.BytesSent),
			formatBytes(iface.BytesRecv),
			fmt.Sprintf("%d", iface.PacketsSent),
			fmt.Sprintf("%d", iface.PacketsRecv),
		})
	}

	table.Render()
}

// PrintProcessInfo 打印进程信息
func PrintProcessInfo(info monitor.ProcessInfo) {
	fmt.Printf("  ")
	colorLabel.Print("进程总数: ")
	colorValue.Println(info.TotalProcesses)

	fmt.Println()
	colorTitle.Println("🔥 Top CPU 进程")
	printProcessTable(info.TopCPU, "cpu")

	fmt.Println()
	colorTitle.Println("💾 Top 内存进程")
	printProcessTable(info.TopMemory, "memory")
}

func printProcessTable(processes []monitor.ProcessDetail, sortBy string) {
	table := tablewriter.NewWriter(os.Stdout)

	if sortBy == "cpu" {
		table.SetHeader([]string{"PID", "用户", "CPU%", "内存", "状态", "命令"})
	} else {
		table.SetHeader([]string{"PID", "用户", "内存", "CPU%", "状态", "命令"})
	}

	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_CENTER,
		tablewriter.ALIGN_LEFT,
	})

	for _, p := range processes {
		cmd := p.Name
		if len(cmd) > 40 {
			cmd = cmd[:37] + "..."
		}

		if sortBy == "cpu" {
			table.Append([]string{
				fmt.Sprintf("%d", p.PID),
				p.Username,
				fmt.Sprintf("%.1f%%", p.CPUPercent),
				fmt.Sprintf("%.1f MB", p.MemoryMB),
				p.Status,
				cmd,
			})
		} else {
			table.Append([]string{
				fmt.Sprintf("%d", p.PID),
				p.Username,
				fmt.Sprintf("%.1f MB", p.MemoryMB),
				fmt.Sprintf("%.1f%%", p.CPUPercent),
				p.Status,
				cmd,
			})
		}
	}

	table.Render()
}

// PrintDockerInfo 打印 Docker 信息（简洁版）
func PrintDockerInfo(info monitor.DockerInfo) {
	colorTitle.Printf("🐳 Docker 容器 (%d 运行中 / %d 总计)\n", info.RunningCount, info.TotalCount)

	if len(info.Containers) == 0 {
		fmt.Printf("  ")
		colorLabel.Println("暂无容器")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"容器名", "镜像", "状态", "CPU", "内存"})
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, c := range info.Containers {
		name := c.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}

		image := c.Image
		if len(image) > 25 {
			image = image[:22] + "..."
		}

		status := c.Status
		if len(status) > 20 {
			status = status[:17] + "..."
		}

		cpuStr := "-"
		memStr := "-"
		if c.State == "running" {
			cpuStr = fmt.Sprintf("%.1f%%", c.CPUPercent)
			memStr = fmt.Sprintf("%.0f MB", c.MemoryUsageMB)
		}

		table.Append([]string{
			name,
			image,
			status,
			cpuStr,
			memStr,
		})
	}

	table.Render()
}

// PrintDockerInfoDetailed 打印 Docker 详细信息
func PrintDockerInfoDetailed(info monitor.DockerInfo) {
	fmt.Printf("  ")
	colorLabel.Print("运行中容器: ")
	colorSuccess.Printf("%d ", info.RunningCount)
	colorLabel.Print("/ 总计: ")
	colorValue.Println(info.TotalCount)

	if len(info.Containers) == 0 {
		fmt.Println()
		fmt.Printf("  ")
		colorLabel.Println("暂无容器")
		return
	}

	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "容器名", "镜像", "端口", "状态", "CPU%", "内存", "运行时长"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, c := range info.Containers {
		name := c.Name
		if len(name) > 15 {
			name = name[:12] + "..."
		}

		image := c.Image
		if len(image) > 18 {
			image = image[:15] + "..."
		}

		// 格式化端口映射
		portStr := formatPorts(c.Ports)
		if portStr == "" {
			portStr = "-"
		}

		cpuStr := "-"
		memStr := "-"

		if c.State == "running" {
			cpuStr = fmt.Sprintf("%.1f%%", c.CPUPercent)
			memStr = fmt.Sprintf("%.0f MB", c.MemoryUsageMB)
		}

		table.Append([]string{
			c.ID,
			name,
			image,
			portStr,
			c.Status,
			cpuStr,
			memStr,
			c.Uptime,
		})
	}

	table.Render()
}

// formatPorts 格式化端口映射
func formatPorts(ports []monitor.PortMapping) string {
	if len(ports) == 0 {
		return ""
	}

	var portStrs []string
	for _, p := range ports {
		if p.PublicPort > 0 {
			portStrs = append(portStrs, fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type))
		} else {
			portStrs = append(portStrs, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}

	// 只显示前2个端口，避免太长
	if len(portStrs) > 2 {
		return strings.Join(portStrs[:2], ", ") + "..."
	}
	return strings.Join(portStrs, ", ")
}

// PrintContainerDetail 打印容器详情
func PrintContainerDetail(info monitor.ContainerInfo) {
	fmt.Printf("  ")
	colorLabel.Print("容器 ID: ")
	colorValue.Println(info.ID)

	fmt.Printf("  ")
	colorLabel.Print("容器名: ")
	colorValue.Println(info.Name)

	fmt.Printf("  ")
	colorLabel.Print("镜像: ")
	colorValue.Println(info.Image)

	fmt.Printf("  ")
	colorLabel.Print("状态: ")
	if info.State == "running" {
		colorSuccess.Println(info.Status)
	} else {
		colorWarning.Println(info.Status)
	}

	fmt.Printf("  ")
	colorLabel.Print("运行时长: ")
	colorInfo.Println(info.Uptime)

	if info.State == "running" {
		fmt.Println()
		colorTitle.Println("📊 资源使用情况")

		fmt.Printf("  ")
		colorLabel.Print("CPU 使用率: ")
		printPercentWithBar(info.CPUPercent, 40)

		fmt.Printf("  ")
		colorLabel.Print("内存使用: ")
		colorValue.Printf("%.1f MB / %.1f MB ", info.MemoryUsageMB, info.MemoryLimitMB)
		printPercentWithBar(info.MemPercent, 30)

		fmt.Printf("  ")
		colorLabel.Print("网络 I/O: ")
		colorSuccess.Printf("↑ %.2f MB  ", info.NetOutputMB)
		colorInfo.Printf("↓ %.2f MB\n", info.NetInputMB)

		fmt.Printf("  ")
		colorLabel.Print("磁盘 I/O: ")
		colorSuccess.Printf("读 %.2f MB  ", info.BlockInputMB)
		colorWarning.Printf("写 %.2f MB\n", info.BlockOutputMB)
	}
}

// 辅助函数

func printPercentWithBar(percent float64, width int) {
	// 颜色选择
	var c *color.Color
	if percent < 50 {
		c = colorSuccess
	} else if percent < 80 {
		c = colorWarning
	} else {
		c = colorError
	}

	// 打印百分比
	c.Printf("%.1f%% ", percent)

	// 打印进度条
	filledWidth := int(percent / 100 * float64(width))
	if filledWidth > width {
		filledWidth = width
	}

	c.Print(strings.Repeat("█", filledWidth))
	colorLabel.Print(strings.Repeat("░", width-filledWidth))
	fmt.Println()
}

func printLoadValue(load float64, cores int) {
	threshold := float64(cores) * 0.7

	if load < threshold {
		colorSuccess.Printf("%.2f\n", load)
	} else if load < float64(cores) {
		colorWarning.Printf("%.2f\n", load)
	} else {
		colorError.Printf("%.2f\n", load)
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%d 天 %d 小时", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%d 小时 %d 分钟", hours, minutes)
	} else {
		return fmt.Sprintf("%d 分钟", minutes)
	}
}

// PrintPortInfo 打印端口信息
func PrintPortInfo(info monitor.PortInfo) {
	if len(info.Listening) == 0 {
		fmt.Printf("  ")
		colorLabel.Println("未检测到监听端口")
		return
	}

	fmt.Printf("  ")
	colorLabel.Print("监听端口总数: ")
	colorValue.Println(len(info.Listening))

	fmt.Println()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"端口", "协议", "地址", "进程", "PID", "状态"})
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, p := range info.Listening {
		processName := p.ProcessName
		if processName == "" {
			processName = "-"
		}

		pidStr := "-"
		if p.PID > 0 {
			pidStr = fmt.Sprintf("%d", p.PID)
		}

		address := p.Address
		if address == "" || address == "0.0.0.0" {
			address = "所有接口"
		}

		table.Append([]string{
			fmt.Sprintf("%d", p.Port),
			p.Protocol,
			address,
			processName,
			pidStr,
			p.State,
		})
	}

	table.Render()
}
