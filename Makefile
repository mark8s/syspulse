.PHONY: build run clean install test fmt

# 构建
build:
	go build -o syspulse main.go

# 运行
run:
	go run main.go

# 清理
clean:
	rm -f syspulse

# 安装到系统
install: build
	sudo cp syspulse /usr/local/bin/

# 测试
test:
	go test -v ./...

# 格式化代码
fmt:
	go fmt ./...

# 下载依赖
deps:
	go mod download
	go mod tidy

# 交叉编译
build-linux:
	GOOS=linux GOARCH=amd64 go build -o syspulse-linux-amd64 main.go

build-all:
	GOOS=linux GOARCH=amd64 go build -o syspulse-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o syspulse-linux-arm64 main.go

