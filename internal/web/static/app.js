// WebSocket 连接
let ws = null;
let reconnectTimer = null;

// 连接 WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?interval=2`;
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = () => {
        console.log('WebSocket 连接成功');
        updateStatus(true);
        clearTimeout(reconnectTimer);
    };
    
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        updateUI(data);
    };
    
    ws.onerror = (error) => {
        console.error('WebSocket 错误:', error);
        updateStatus(false);
    };
    
    ws.onclose = () => {
        console.log('WebSocket 连接关闭，3秒后重连...');
        updateStatus(false);
        reconnectTimer = setTimeout(connectWebSocket, 3000);
    };
}

// 更新连接状态
function updateStatus(connected) {
    const statusEl = document.getElementById('status');
    const lastUpdateEl = document.getElementById('last-update');
    
    if (connected) {
        statusEl.classList.remove('disconnected');
        lastUpdateEl.textContent = '已连接';
    } else {
        statusEl.classList.add('disconnected');
        lastUpdateEl.textContent = '连接断开';
    }
}

// 更新界面
function updateUI(data) {
    document.getElementById('last-update').textContent = `最后更新: ${new Date().toLocaleTimeString()}`;
    
    // 系统信息
    if (data.system) {
        document.getElementById('hostname').textContent = data.system.Hostname || '-';
        document.getElementById('os').textContent = data.system.OS || '-';
        document.getElementById('kernel').textContent = data.system.Kernel || '-';
        document.getElementById('uptime').textContent = formatUptime(data.system.Uptime) || '-';
    }
    
    // CPU
    if (data.cpu) {
        const cpuPercent = data.cpu.UsagePercent.toFixed(1);
        document.getElementById('cpu-percent').textContent = cpuPercent + '%';
        updateProgressBar('cpu-bar', cpuPercent);
        document.getElementById('cpu-cores').textContent = data.cpu.CoreCount;
        document.getElementById('cpu-load').textContent = 
            `${data.cpu.LoadAvg1.toFixed(2)} / ${data.cpu.LoadAvg5.toFixed(2)} / ${data.cpu.LoadAvg15.toFixed(2)}`;
    }
    
    // 内存
    if (data.memory) {
        const memPercent = data.memory.UsedPercent.toFixed(1);
        const memUsed = formatBytes(data.memory.Used);
        const memTotal = formatBytes(data.memory.Total);
        document.getElementById('mem-usage').textContent = `${memUsed} / ${memTotal}`;
        updateProgressBar('mem-bar', memPercent);
        
        const swapPercent = data.memory.SwapPercent.toFixed(1);
        const swapUsed = formatBytes(data.memory.SwapUsed);
        const swapTotal = formatBytes(data.memory.SwapTotal);
        document.getElementById('swap-usage').textContent = `${swapUsed} / ${swapTotal}`;
        updateProgressBar('swap-bar', swapPercent);
    }
    
    // 磁盘
    if (data.disk && data.disk.Partitions) {
        updateDiskList(data.disk.Partitions);
    }
    
    // 网络
    if (data.network && data.network.Interfaces) {
        updateNetworkList(data.network.Interfaces);
    }
    
    // 端口
    if (data.ports && data.ports.Listening) {
        updatePortList(data.ports.Listening);
    }
    
    // Docker
    if (data.docker) {
        updateDockerList(data.docker);
    }
    
    // 进程
    if (data.process && data.process.TopCPU) {
        updateProcessTable(data.process.TopCPU);
    }
}

// 更新进度条
function updateProgressBar(id, percent) {
    const bar = document.getElementById(id);
    bar.style.width = percent + '%';
    
    // 根据使用率改变颜色
    bar.classList.remove('warning', 'danger');
    if (percent > 80) {
        bar.classList.add('danger');
    } else if (percent > 50) {
        bar.classList.add('warning');
    }
}

// 更新磁盘列表（df -h 风格）
function updateDiskList(partitions) {
    const container = document.getElementById('disk-list');
    container.innerHTML = '';
    
    // 创建表格
    const table = document.createElement('table');
    table.innerHTML = `
        <thead>
            <tr>
                <th>文件系统</th>
                <th>容量</th>
                <th>已用</th>
                <th>可用</th>
                <th>已用%</th>
                <th>挂载点</th>
            </tr>
        </thead>
        <tbody>
            ${partitions.map(partition => {
                const percent = partition.UsedPercent.toFixed(0);
                const percentClass = getPercentClass(partition.UsedPercent);
                const used = formatBytes(partition.Used);
                const total = formatBytes(partition.Total);
                const free = formatBytes(partition.Free);
                
                return `
                    <tr>
                        <td class="breakable"><code>${partition.Device}</code></td>
                        <td class="nowrap">${total}</td>
                        <td class="nowrap">${used}</td>
                        <td class="nowrap">${free}</td>
                        <td class="nowrap">
                            <span class="percent-badge ${percentClass}">${percent}%</span>
                            <div class="progress-bar-mini">
                                <div class="progress-fill ${percentClass}" style="width: ${percent}%"></div>
                            </div>
                        </td>
                        <td class="breakable"><strong>${partition.Mountpoint}</strong></td>
                    </tr>
                `;
            }).join('')}
        </tbody>
    `;
    
    container.appendChild(table);
}

// 更新网络列表
function updateNetworkList(interfaces) {
    const container = document.getElementById('network-list');
    container.innerHTML = '';
    
    interfaces.forEach(iface => {
        if (iface.Addrs && iface.Addrs.length > 0) {
            const div = document.createElement('div');
            div.className = 'network-item';
            div.innerHTML = `
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                    <div>
                        <strong style="color: var(--primary); font-size: 1.1em;">${iface.Name}</strong>
                        <span style="margin-left: 15px; color: var(--text);">${iface.Addrs[0]}</span>
                    </div>
                </div>
                <div class="network-stats">
                    <span><strong>📤 发送:</strong> ${formatBytes(iface.BytesSent)}</span>
                    <span><strong>📥 接收:</strong> ${formatBytes(iface.BytesRecv)}</span>
                    <span><strong>📦 发送包:</strong> ${iface.PacketsSent.toLocaleString()}</span>
                    <span><strong>📦 接收包:</strong> ${iface.PacketsRecv.toLocaleString()}</span>
                </div>
            `;
            container.appendChild(div);
        }
    });
}

// 更新端口列表
function updatePortList(ports) {
    const container = document.getElementById('port-list');
    if (!container) return;
    
    container.innerHTML = '';
    
    if (!ports || ports.length === 0) {
        container.innerHTML = '<div style="text-align: center; color: var(--text-muted); padding: 20px;">未检测到监听端口</div>';
        return;
    }
    
    // 按端口号排序
    ports.sort((a, b) => a.Port - b.Port);
    
    // 创建表格
    const table = document.createElement('table');
    table.innerHTML = `
        <thead>
            <tr>
                <th>端口</th>
                <th>协议</th>
                <th>地址</th>
                <th>进程</th>
                <th>状态</th>
            </tr>
        </thead>
        <tbody>
            ${ports.map(port => `
                <tr>
                    <td class="nowrap"><strong>${port.Port}</strong></td>
                    <td class="nowrap">${port.Protocol.toUpperCase()}</td>
                    <td class="nowrap">${port.Address || '0.0.0.0'}</td>
                    <td class="nowrap">${port.ProcessName || '-'} ${port.PID ? `(${port.PID})` : ''}</td>
                    <td class="nowrap"><span class="status-badge status-${port.State.toLowerCase()}">${port.State}</span></td>
                </tr>
            `).join('')}
        </tbody>
    `;
    container.appendChild(table);
}

// 更新 Docker 列表
function updateDockerList(docker) {
    const statusEl = document.getElementById('docker-status');
    const container = document.getElementById('docker-list');
    
    if (!docker.Available) {
        statusEl.innerHTML = '<span style="color: var(--danger);">❌ Docker 不可用</span>';
        container.innerHTML = '<div style="text-align: center; color: var(--text-muted); padding: 20px;">Docker 服务未运行</div>';
        return;
    }
    
    statusEl.innerHTML = `<span style="color: var(--success);">✅ 运行中: ${docker.RunningCount} / 总计: ${docker.TotalCount}</span>`;
    
    if (!docker.Containers || docker.Containers.length === 0) {
        container.innerHTML = '<div style="text-align: center; color: var(--text-muted); padding: 20px;">暂无容器</div>';
        return;
    }
    
    // 创建表格（显示完整内容，不省略）
    const table = document.createElement('table');
    table.innerHTML = `
        <thead>
            <tr>
                <th>容器名</th>
                <th>镜像</th>
                <th>端口</th>
                <th>状态</th>
                <th>CPU</th>
                <th>内存</th>
            </tr>
        </thead>
        <tbody>
            ${docker.Containers.map(c => {
                const isRunning = c.State === 'running';
                const statusClass = isRunning ? 'status-running' : 'status-stopped';
                const cpu = isRunning ? `${c.CPUPercent.toFixed(1)}%` : '-';
                const mem = isRunning ? `${c.MemoryUsageMB.toFixed(0)} MB` : '-';
                
                // 格式化端口（显示所有端口）
                let ports = '-';
                if (c.Ports && c.Ports.length > 0) {
                    const portList = c.Ports.map(p => {
                        if (p.PublicPort > 0) {
                            return `${p.PublicPort}→${p.PrivatePort}`;
                        }
                        return `${p.PrivatePort}`;
                    });
                    ports = portList.join(', ');
                }
                
                return `
                    <tr>
                        <td class="nowrap"><strong>${c.Name}</strong></td>
                        <td class="breakable">${c.Image}</td>
                        <td class="nowrap"><span class="port-badge">${ports}</span></td>
                        <td class="nowrap"><span class="${statusClass}">${c.Status}</span></td>
                        <td>${cpu}</td>
                        <td>${mem}</td>
                    </tr>
                `;
            }).join('')}
        </tbody>
    `;
    
    container.innerHTML = '';
    container.appendChild(table);
}

// 更新进程表格
function updateProcessTable(processes) {
    const tbody = document.getElementById('process-tbody');
    tbody.innerHTML = '';
    
    processes.forEach(p => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${p.PID}</td>
            <td class="nowrap">${p.Username}</td>
            <td><strong>${p.CPUPercent.toFixed(1)}%</strong></td>
            <td>${p.MemoryMB.toFixed(1)} MB</td>
            <td class="breakable" title="${p.Name}">${p.Name}</td>
        `;
        tbody.appendChild(tr);
    });
}

// 工具函数
function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i];
}

function formatUptime(seconds) {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    
    if (days > 0) {
        return `${days} 天 ${hours} 小时`;
    } else if (hours > 0) {
        return `${hours} 小时 ${minutes} 分钟`;
    } else {
        return `${minutes} 分钟`;
    }
}

function getPercentClass(percent) {
    if (percent > 80) return 'danger';
    if (percent > 50) return 'warning';
    return '';
}

function truncate(str, len) {
    if (str.length <= len) return str;
    return str.substring(0, len - 3) + '...';
}

// 页面加载时连接
window.addEventListener('load', () => {
    connectWebSocket();
});

// 页面卸载时关闭连接
window.addEventListener('beforeunload', () => {
    if (ws) {
        ws.close();
    }
});

