#!/bin/bash

# SysPulse 安装脚本

set -e

echo "🚀 正在安装 SysPulse..."
echo ""

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到 Go，请先安装 Go 1.21 或更高版本"
    echo "   访问: https://golang.org/dl/"
    exit 1
fi

# 检查 Go 版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ 错误: Go 版本过低，需要 $REQUIRED_VERSION 或更高版本"
    echo "   当前版本: $GO_VERSION"
    exit 1
fi

echo "✅ Go 版本: $GO_VERSION"
echo ""

# 下载依赖
echo "📦 正在下载依赖..."
go mod download
go mod tidy
echo ""

# 构建
echo "🔨 正在构建..."
go build -o syspulse main.go
echo ""

# 安装到系统
if [ "$EUID" -eq 0 ]; then
    echo "📦 正在安装到 /usr/local/bin/..."
    cp syspulse /usr/local/bin/
    chmod +x /usr/local/bin/syspulse
    echo ""
    echo "✅ 安装成功！"
    echo ""
    echo "现在你可以在任何地方使用 'syspulse' 命令了！"
else
    echo "✅ 构建成功！"
    echo ""
    echo "可执行文件位于: ./syspulse"
    echo ""
    echo "如需安装到系统路径，请运行:"
    echo "  sudo make install"
    echo ""
    echo "或者直接运行:"
    echo "  ./syspulse"
fi

echo ""
echo "📖 快速开始:"
echo "  syspulse              # 显示仪表盘"
echo "  syspulse docker       # 查看 Docker 容器"
echo "  syspulse cpu          # 查看 CPU 信息"
echo "  syspulse --help       # 查看所有命令"
echo ""
echo "🎉 享受使用吧！"

