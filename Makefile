.PHONY: all build run clean help swag test docker

BINARY=gin_pipeline
MAIN_GO=main.go

all: build

build:
	go build -o $(BINARY) $(MAIN_GO)

run:
	go run $(MAIN_GO)

clean:
	go clean
	rm -f $(BINARY)

swag:
	swag init

test:
	go test -v ./...

docker:
	docker build -t $(BINARY) .

docker-run:
	docker run -p 8080:8080 $(BINARY)

help:
	@echo "make - 编译项目"
	@echo "make build - 编译项目"
	@echo "make run - 运行项目"
	@echo "make clean - 清理编译文件"
	@echo "make swag - 生成 Swagger 文档"
	@echo "make test - 运行测试"
	@echo "make docker - 构建 Docker 镜像"
	@echo "make docker-run - 运行 Docker 容器"
