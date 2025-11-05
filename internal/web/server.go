package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//go:embed static/*
var staticFiles embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Server Web 服务器
type Server struct {
	host   string
	port   int
	router *mux.Router
}

// NewServer 创建新的 Web 服务器
func NewServer(host string, port int) *Server {
	s := &Server{
		host:   host,
		port:   port,
		router: mux.NewRouter(),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// API 路由
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/system", handleSystem).Methods("GET")
	api.HandleFunc("/cpu", handleCPU).Methods("GET")
	api.HandleFunc("/memory", handleMemory).Methods("GET")
	api.HandleFunc("/disk", handleDisk).Methods("GET")
	api.HandleFunc("/network", handleNetwork).Methods("GET")
	api.HandleFunc("/port", handlePort).Methods("GET")
	api.HandleFunc("/process", handleProcess).Methods("GET")
	api.HandleFunc("/docker", handleDocker).Methods("GET")
	api.HandleFunc("/docker/{id}", handleDockerDetail).Methods("GET")
	api.HandleFunc("/all", handleAll).Methods("GET")

	// WebSocket 路由
	s.router.HandleFunc("/ws", handleWebSocket)

	// 静态文件服务
	staticFS, _ := fs.Sub(staticFiles, "static")
	s.router.PathPrefix("/").Handler(http.FileServer(http.FS(staticFS)))
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv.ListenAndServe()
}
