# SysPulse 使用示例

## 基础用法

### 1. 查看系统仪表盘（推荐）

显示所有系统资源的概览：

```bash
syspulse
# 或
syspulse dashboard
```

### 2. 实时监控模式

每 2 秒自动刷新：

```bash
syspulse dashboard --watch
```

自定义刷新间隔（5 秒）：

```bash
syspulse dashboard --watch --interval 5
```

## 查看特定资源

### CPU 信息

```bash
syspulse cpu
```

显示内容：
- CPU 型号
- 总体使用率（带进度条）
- 每个核心的使用率
- 负载平均值（1分钟/5分钟/15分钟）

### 内存信息

```bash
syspulse memory
```

显示内容：
- 物理内存：总量、已用、可用、缓存、缓冲
- Swap 交换分区使用情况
- 使用率可视化

### 磁盘信息

```bash
syspulse disk
```

显示内容：
- 所有磁盘分区列表
- 挂载点、设备、文件系统类型
- 容量、已用、可用空间
- 使用率百分比

### 网络信息

```bash
syspulse network
```

显示内容：
- 网络接口列表
- IP 地址
- 发送/接收字节数和包数
- 流量统计

### 进程信息

查看 Top 10 进程：

```bash
syspulse process
```

查看 Top 20 进程：

```bash
syspulse process --top 20
```

显示内容：
- Top N CPU 占用进程
- Top N 内存占用进程
- 进程详细信息（PID、用户、CPU、内存、状态）

## Docker 容器监控

### 查看所有容器

```bash
syspulse docker
```

显示内容：
- 容器列表（运行中/全部）
- 容器名、镜像、状态
- CPU 和内存使用情况
- 运行时长

### 实时监控容器

```bash
syspulse docker --watch
```

自定义刷新间隔：

```bash
syspulse docker --watch --interval 3
```

### 查看特定容器详情

使用容器 ID（完整或短 ID）：

```bash
syspulse docker --container abc123def456
# 或使用短 ID
syspulse docker --container abc123
```

显示内容：
- 容器基本信息
- CPU 使用率（带进度条）
- 内存使用情况
- 网络 I/O（上传/下载）
- 磁盘 I/O（读/写）

## 实用技巧

### 1. 组合使用 watch 命令

如果想要更灵活的刷新控制：

```bash
watch -n 1 syspulse cpu
```

### 2. 输出到文件

```bash
syspulse > system_report.txt
```

### 3. 只查看 Docker 运行中的容器

配合 grep 使用：

```bash
syspulse docker | grep "运行中"
```

### 4. 监控特定磁盘

```bash
syspulse disk | grep "/home"
```

### 5. 在 SSH 会话中使用

```bash
ssh user@server "syspulse"
```

### 6. 定时检查并记录

使用 crontab 定时记录系统状态：

```bash
# 每小时记录一次
0 * * * * /usr/local/bin/syspulse >> /var/log/syspulse.log 2>&1
```

## 常见使用场景

### 场景 1: 排查系统性能问题

```bash
# 1. 先看总体情况
syspulse

# 2. 发现 CPU 高，查看详情
syspulse cpu

# 3. 看看是哪些进程
syspulse process --top 20

# 4. 检查是否是容器引起的
syspulse docker
```

### 场景 2: 监控 Docker 容器

```bash
# 实时监控所有容器
syspulse docker --watch

# 发现某个容器资源占用高，查看详情
syspulse docker --container <container-id>
```

### 场景 3: 检查磁盘空间

```bash
# 查看所有分区
syspulse disk

# 如果某个分区快满了，可以进一步排查
du -sh /* | sort -hr | head -20
```

### 场景 4: 持续监控系统

```bash
# 在 tmux 或 screen 中运行
tmux new -s monitor
syspulse dashboard --watch --interval 2
# 按 Ctrl+B 然后 D 分离会话
```

## 命令速查表

| 命令 | 说明 |
|------|------|
| `syspulse` | 显示仪表盘 |
| `syspulse dashboard --watch` | 实时仪表盘 |
| `syspulse cpu` | CPU 信息 |
| `syspulse memory` | 内存信息 |
| `syspulse disk` | 磁盘信息 |
| `syspulse network` | 网络信息 |
| `syspulse process` | 进程信息 |
| `syspulse docker` | Docker 容器 |
| `syspulse docker --watch` | 实时监控容器 |
| `syspulse --help` | 帮助信息 |

## 故障排除

### Docker 不可用

如果看到 "Docker 不可用" 消息：

1. 检查 Docker 是否安装：
   ```bash
   docker --version
   ```

2. 检查 Docker 服务是否运行：
   ```bash
   sudo systemctl status docker
   ```

3. 确保当前用户有 Docker 权限：
   ```bash
   sudo usermod -aG docker $USER
   # 重新登录后生效
   ```

### 权限问题

某些信息可能需要 root 权限：

```bash
sudo syspulse
```

### 命令不存在

如果显示 "command not found"：

```bash
# 使用完整路径
./syspulse

# 或添加到 PATH
export PATH=$PATH:$(pwd)
```

## 进阶使用

### 与其他工具集成

#### 与 Prometheus 集成（需要自定义脚本）

```bash
# 导出为 JSON 格式（需要添加 --json 参数支持）
syspulse --json | curl -X POST http://prometheus-pushgateway:9091/metrics/job/syspulse -d @-
```

#### 发送告警通知

```bash
# 检查 CPU 使用率，超过 80% 发送通知
CPU_USAGE=$(syspulse cpu | grep "总体使用率" | awk '{print $3}' | sed 's/%//')
if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
    echo "CPU 使用率过高: $CPU_USAGE%" | mail -s "Alert" admin@example.com
fi
```

## 反馈和贡献

如果你有任何问题或建议，欢迎提交 Issue 或 Pull Request！

