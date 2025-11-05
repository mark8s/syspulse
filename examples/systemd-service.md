# 将 SysPulse 作为 systemd 服务运行

如果你想定期记录系统状态，可以将 SysPulse 配置为 systemd 服务。

## 方案 1: 使用 systemd Timer（推荐）

### 1. 创建服务文件

创建 `/etc/systemd/system/syspulse-monitor.service`:

```ini
[Unit]
Description=SysPulse System Monitor
After=network.target

[Service]
Type=oneshot
User=root
ExecStart=/usr/local/bin/syspulse
StandardOutput=append:/var/log/syspulse/monitor.log
StandardError=append:/var/log/syspulse/error.log

[Install]
WantedBy=multi-user.target
```

### 2. 创建 Timer 文件

创建 `/etc/systemd/system/syspulse-monitor.timer`:

```ini
[Unit]
Description=Run SysPulse Monitor every 5 minutes
Requires=syspulse-monitor.service

[Timer]
OnBootSec=1min
OnUnitActiveSec=5min
AccuracySec=1s

[Install]
WantedBy=timers.target
```

### 3. 创建日志目录

```bash
sudo mkdir -p /var/log/syspulse
sudo chmod 755 /var/log/syspulse
```

### 4. 启用并启动

```bash
# 重载 systemd
sudo systemctl daemon-reload

# 启用定时器
sudo systemctl enable syspulse-monitor.timer

# 启动定时器
sudo systemctl start syspulse-monitor.timer

# 查看状态
sudo systemctl status syspulse-monitor.timer

# 查看日志
sudo journalctl -u syspulse-monitor.service -f
```

## 方案 2: Docker 容器监控服务

专门用于监控 Docker 容器的服务：

### 1. 创建服务文件

创建 `/etc/systemd/system/syspulse-docker.service`:

```ini
[Unit]
Description=SysPulse Docker Monitor
After=docker.service
Requires=docker.service

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/syspulse docker --watch --interval 10
StandardOutput=append:/var/log/syspulse/docker.log
StandardError=append:/var/log/syspulse/docker-error.log
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 2. 启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl enable syspulse-docker.service
sudo systemctl start syspulse-docker.service
sudo systemctl status syspulse-docker.service
```

## 日志轮转配置

创建 `/etc/logrotate.d/syspulse`:

```
/var/log/syspulse/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 root root
    sharedscripts
    postrotate
        systemctl reload syspulse-monitor.timer > /dev/null 2>&1 || true
    endscript
}
```

## 方案 3: 告警监控服务

监控系统并在异常时发送告警：

### 1. 创建监控脚本

创建 `/usr/local/bin/syspulse-alert.sh`:

```bash
#!/bin/bash

LOG_FILE="/var/log/syspulse/alert.log"
ALERT_EMAIL="admin@example.com"

# 获取 CPU 使用率
CPU_USAGE=$(syspulse cpu | grep -oP '\d+\.\d+(?=%)')

# 获取内存使用率
MEM_USAGE=$(syspulse memory | grep -oP '\d+\.\d+(?=%)')

# 获取磁盘使用率
DISK_USAGE=$(syspulse disk | grep -oP '\d+\.\d+(?=%)' | head -1)

# 检查阈值
if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
    echo "[$(date)] CPU 使用率过高: $CPU_USAGE%" >> $LOG_FILE
    echo "警告: CPU 使用率 $CPU_USAGE%" | mail -s "SysPulse Alert: High CPU" $ALERT_EMAIL
fi

if (( $(echo "$MEM_USAGE > 90" | bc -l) )); then
    echo "[$(date)] 内存使用率过高: $MEM_USAGE%" >> $LOG_FILE
    echo "警告: 内存使用率 $MEM_USAGE%" | mail -s "SysPulse Alert: High Memory" $ALERT_EMAIL
fi

if (( $(echo "$DISK_USAGE > 85" | bc -l) )); then
    echo "[$(date)] 磁盘使用率过高: $DISK_USAGE%" >> $LOG_FILE
    echo "警告: 磁盘使用率 $DISK_USAGE%" | mail -s "SysPulse Alert: High Disk" $ALERT_EMAIL
fi
```

### 2. 赋予执行权限

```bash
sudo chmod +x /usr/local/bin/syspulse-alert.sh
```

### 3. 创建服务和定时器

类似方案 1，但 ExecStart 指向 `/usr/local/bin/syspulse-alert.sh`

## 查看和管理

```bash
# 查看所有 syspulse 相关服务
systemctl list-units | grep syspulse

# 查看定时器列表
systemctl list-timers

# 停止服务
sudo systemctl stop syspulse-monitor.timer

# 禁用服务
sudo systemctl disable syspulse-monitor.timer

# 查看日志
sudo tail -f /var/log/syspulse/monitor.log

# 查看最近的运行记录
sudo journalctl -u syspulse-monitor.service --since today
```

## 卸载

```bash
# 停止并禁用
sudo systemctl stop syspulse-monitor.timer
sudo systemctl disable syspulse-monitor.timer

# 删除文件
sudo rm /etc/systemd/system/syspulse-monitor.service
sudo rm /etc/systemd/system/syspulse-monitor.timer
sudo rm /etc/logrotate.d/syspulse
sudo rm -rf /var/log/syspulse

# 重载 systemd
sudo systemctl daemon-reload
```

