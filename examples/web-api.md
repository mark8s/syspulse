# SysPulse Web API 文档

## 概述

SysPulse 提供了完整的 Web 界面和 RESTful API，可以通过 HTTP 和 WebSocket 获取系统监控数据。

## 启动 Web 服务器

```bash
# 默认端口 3000
./syspulse web

# 自定义端口
./syspulse web --port 8888

# 指定监听地址（只允许本地访问）
./syspulse web --host 127.0.0.1

# 允许所有网络访问
./syspulse web --host 0.0.0.0 --port 3000
```

## RESTful API

### 基础 URL

```
http://localhost:3000/api
```

### 端点列表

#### 1. 获取系统信息

```http
GET /api/system
```

**响应示例：**
```json
{
  "Hostname": "server-01",
  "OS": "linux",
  "Kernel": "5.15.0-91-generic",
  "Uptime": 1296000,
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 2. 获取 CPU 信息

```http
GET /api/cpu
```

**响应示例：**
```json
{
  "UsagePercent": 45.3,
  "CoreCount": 8,
  "ModelName": "Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz",
  "LoadAvg1": 2.1,
  "LoadAvg5": 2.5,
  "LoadAvg15": 2.8,
  "PerCoreUsage": [52.1, 43.7, 38.2, 61.5, 45.0, 42.3, 47.8, 39.6],
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 3. 获取内存信息

```http
GET /api/memory
```

**响应示例：**
```json
{
  "Total": 17179869184,
  "Used": 13421772800,
  "Available": 3758096384,
  "UsedPercent": 78.1,
  "SwapTotal": 8589934592,
  "SwapUsed": 2147483648,
  "SwapPercent": 25.0,
  "Cached": 4294967296,
  "Buffers": 536870912,
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 4. 获取磁盘信息

```http
GET /api/disk
```

**响应示例：**
```json
{
  "Partitions": [
    {
      "Device": "/dev/sda1",
      "Mountpoint": "/",
      "Fstype": "ext4",
      "Total": 214748364800,
      "Used": 91625968640,
      "Free": 123122396160,
      "UsedPercent": 42.7
    }
  ],
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 5. 获取网络信息

```http
GET /api/network
```

**响应示例：**
```json
{
  "Interfaces": [
    {
      "Name": "eth0",
      "BytesSent": 2684354560,
      "BytesRecv": 16442450944,
      "PacketsSent": 1234567,
      "PacketsRecv": 7654321,
      "Addrs": ["192.168.1.100/24"]
    }
  ],
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 6. 获取端口监听信息

```http
GET /api/port
```

**响应示例：**
```json
{
  "Listening": [
    {
      "Port": 22,
      "Address": "0.0.0.0",
      "Protocol": "tcp",
      "State": "LISTEN",
      "PID": 1234,
      "ProcessName": "sshd"
    },
    {
      "Port": 80,
      "Address": "0.0.0.0",
      "Protocol": "tcp",
      "State": "LISTEN",
      "PID": 5678,
      "ProcessName": "nginx"
    }
  ],
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 7. 获取进程信息

```http
GET /api/process?top=10
```

**查询参数：**
- `top` - 返回 Top N 进程（默认 10）

**响应示例：**
```json
{
  "TotalProcesses": 156,
  "TopCPU": [
    {
      "PID": 1234,
      "Name": "chrome",
      "Username": "user",
      "CPUPercent": 15.7,
      "MemoryMB": 1234.5,
      "MemPercent": 7.6,
      "Status": "R",
      "Command": "/usr/bin/chrome"
    }
  ],
  "TopMemory": [...],
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 8. 获取 Docker 信息

```http
GET /api/docker
```

**响应示例：**
```json
{
  "Available": true,
  "RunningCount": 3,
  "TotalCount": 5,
  "Containers": [
    {
      "ID": "abc123def456",
      "Name": "nginx-web",
      "Image": "nginx:latest",
      "Status": "Up 5 hours",
      "State": "running",
      "CPUPercent": 2.3,
      "MemoryUsageMB": 128.5,
      "MemoryLimitMB": 2048.0,
      "MemPercent": 6.3,
      "NetInputMB": 15.3,
      "NetOutputMB": 2.5,
      "BlockInputMB": 0.5,
      "BlockOutputMB": 1.2,
      "Created": "2025-11-01T10:00:00Z",
      "Uptime": "5h"
    }
  ],
  "Timestamp": "2025-11-05T10:30:00Z"
}
```

#### 9. 获取特定容器详情

```http
GET /api/docker/{container_id}
```

#### 10. 获取所有信息

```http
GET /api/all
```

返回包含所有模块数据的综合响应。

## WebSocket API

### 连接端点

```
ws://localhost:3000/ws
```

### 查询参数

- `interval` - 数据推送间隔（秒），默认 2

### 连接示例

```javascript
const ws = new WebSocket('ws://localhost:3000/ws?interval=2');

ws.onopen = () => {
  console.log('已连接');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('收到数据:', data);
  // data 包含所有监控数据：system, cpu, memory, disk, network, ports, docker, process
};

ws.onerror = (error) => {
  console.error('错误:', error);
};

ws.onclose = () => {
  console.log('连接关闭');
};
```

### 数据格式

WebSocket 推送的数据格式与 `/api/all` 返回的格式相同，包含所有监控模块的实时数据。

## 使用示例

### cURL

```bash
# 获取 CPU 信息
curl http://localhost:3000/api/cpu

# 获取所有信息
curl http://localhost:3000/api/all

# 获取 Top 20 进程
curl http://localhost:3000/api/process?top=20
```

### Python

```python
import requests
import json

# 获取系统信息
response = requests.get('http://localhost:3000/api/system')
data = response.json()
print(json.dumps(data, indent=2))

# 获取所有信息
response = requests.get('http://localhost:3000/api/all')
all_data = response.json()
print(f"CPU 使用率: {all_data['cpu']['UsagePercent']}%")
print(f"内存使用率: {all_data['memory']['UsedPercent']}%")
```

### JavaScript (浏览器)

```javascript
// 获取 Docker 信息
fetch('http://localhost:3000/api/docker')
  .then(response => response.json())
  .then(data => {
    console.log('Docker 容器数:', data.TotalCount);
    console.log('运行中:', data.RunningCount);
  });

// 实时监控 (WebSocket)
const ws = new WebSocket('ws://localhost:3000/ws?interval=1');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  updateDashboard(data);
};

function updateDashboard(data) {
  document.getElementById('cpu').textContent = data.cpu.UsagePercent + '%';
  document.getElementById('memory').textContent = data.memory.UsedPercent + '%';
}
```

## CORS 支持

API 默认启用 CORS，允许跨域访问：

```
Access-Control-Allow-Origin: *
```

## 集成到监控系统

### Prometheus

可以编写 exporter 来采集数据：

```python
from prometheus_client import start_http_server, Gauge
import requests
import time

cpu_usage = Gauge('syspulse_cpu_usage', 'CPU usage percentage')
memory_usage = Gauge('syspulse_memory_usage', 'Memory usage percentage')

def collect_metrics():
    while True:
        data = requests.get('http://localhost:3000/api/all').json()
        cpu_usage.set(data['cpu']['UsagePercent'])
        memory_usage.set(data['memory']['UsedPercent'])
        time.sleep(10)

if __name__ == '__main__':
    start_http_server(9100)
    collect_metrics()
```

### Grafana

1. 配置数据源为 Prometheus
2. 使用上面的 exporter
3. 创建仪表盘展示数据

## 安全建议

1. **生产环境**：使用 `--host 127.0.0.1` 限制只能本地访问
2. **反向代理**：通过 Nginx 添加认证和 SSL
3. **防火墙**：限制访问 IP 范围

### Nginx 反向代理示例

```nginx
server {
    listen 443 ssl;
    server_name monitor.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    auth_basic "Restricted";
    auth_basic_user_file /etc/nginx/.htpasswd;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

## 错误处理

API 使用标准 HTTP 状态码：

- `200 OK` - 成功
- `400 Bad Request` - 请求参数错误
- `404 Not Found` - 资源不存在
- `500 Internal Server Error` - 服务器错误

## 限制

- WebSocket 连接数限制：建议不超过 100 个并发连接
- API 请求频率：无限制，但建议不要过于频繁（推荐 >= 1秒间隔）
- 数据保留：实时数据，不保留历史记录

## 开发

想要扩展 API？查看 `internal/web/handlers.go` 添加新的端点。

