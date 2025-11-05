package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"syspulse/internal/monitor"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// handleSystem 处理系统信息请求
func handleSystem(w http.ResponseWriter, r *http.Request) {
	info := monitor.GetSystemInfo()
	respondJSON(w, info)
}

// handleCPU 处理 CPU 信息请求
func handleCPU(w http.ResponseWriter, r *http.Request) {
	info := monitor.GetCPUInfo()
	respondJSON(w, info)
}

// handleMemory 处理内存信息请求
func handleMemory(w http.ResponseWriter, r *http.Request) {
	info := monitor.GetMemoryInfo()
	respondJSON(w, info)
}

// handleDisk 处理磁盘信息请求
func handleDisk(w http.ResponseWriter, r *http.Request) {
	info := monitor.GetDiskInfo()
	respondJSON(w, info)
}

// handleNetwork 处理网络信息请求
func handleNetwork(w http.ResponseWriter, r *http.Request) {
	info := monitor.GetNetworkInfo()
	respondJSON(w, info)
}

// handleProcess 处理进程信息请求
func handleProcess(w http.ResponseWriter, r *http.Request) {
	topN := 10
	if topNStr := r.URL.Query().Get("top"); topNStr != "" {
		if n, err := strconv.Atoi(topNStr); err == nil {
			topN = n
		}
	}
	info := monitor.GetProcessInfo(topN)
	respondJSON(w, info)
}

// handleDocker 处理 Docker 信息请求
func handleDocker(w http.ResponseWriter, r *http.Request) {
	info := monitor.GetDockerInfo()
	respondJSON(w, info)
}

// handleDockerDetail 处理 Docker 容器详情请求
func handleDockerDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["id"]
	info := monitor.GetContainerDetail(containerID)
	respondJSON(w, info)
}

// handleAll 处理所有信息请求
func handleAll(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"system":  monitor.GetSystemInfo(),
		"cpu":     monitor.GetCPUInfo(),
		"memory":  monitor.GetMemoryInfo(),
		"disk":    monitor.GetDiskInfo(),
		"network": monitor.GetNetworkInfo(),
		"ports":   monitor.GetPortInfo(),
		"docker":  monitor.GetDockerInfo(),
	}
	respondJSON(w, data)
}

// handlePort 处理端口信息请求
func handlePort(w http.ResponseWriter, r *http.Request) {
	info := monitor.GetPortInfo()
	respondJSON(w, info)
}

// handleWebSocket 处理 WebSocket 连接
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 获取刷新间隔（默认 2 秒）
	interval := 2 * time.Second
	if intervalStr := r.URL.Query().Get("interval"); intervalStr != "" {
		if seconds, err := strconv.Atoi(intervalStr); err == nil {
			interval = time.Duration(seconds) * time.Second
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 立即发送第一次数据
	sendAllData(conn)

	// 定期发送更新
	for range ticker.C {
		if err := sendAllData(conn); err != nil {
			break
		}
	}
}

func sendAllData(conn *websocket.Conn) error {
	data := map[string]interface{}{
		"system":  monitor.GetSystemInfo(),
		"cpu":     monitor.GetCPUInfo(),
		"memory":  monitor.GetMemoryInfo(),
		"disk":    monitor.GetDiskInfo(),
		"network": monitor.GetNetworkInfo(),
		"ports":   monitor.GetPortInfo(),
		"docker":  monitor.GetDockerInfo(),
		"process": monitor.GetProcessInfo(10),
	}

	return conn.WriteJSON(data)
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}
