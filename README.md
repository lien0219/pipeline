# 打造企业级标准化CICD平台

# 后端 (Go)

1. **配置管理**：使用YAML配置文件，支持不同环境配置
2. **日志系统**：集成zap日志库，支持日志分割和级别控制
3. **数据库集成**：GORM + MySQL，包含连接池配置
4. **Redis集成**：支持缓存和会话管理
5. **认证系统**：JWT认证中间件
6. **API文档**：集成Swagger (swag)
7. **中间件**：日志、恢复、CORS等
8. **路由管理**：分组和版本控制
9. **错误处理**：统一的错误响应格式


# 前端 (Vue 3)

1. **项目结构**：符合Vue 3最佳实践
2. **状态管理**：使用Pinia
3. **路由**：Vue Router配置，包含路由守卫
4. **API层**：Axios封装，拦截器配置
5. **认证**：完整的登录/注册/个人资料功能
6. **布局**：响应式布局组件
7. **页面**：首页、仪表盘、个人资料等

# Web框架
go get -u github.com/gin-gonic/gin

# ORM和数据库驱动
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# Redis客户端
go get -u github.com/go-redis/redis/v8

# JWT认证
go get -u github.com/dgrijalva/jwt-go

# 日志库
go get -u go.uber.org/zap
go get -u gopkg.in/natefinch/lumberjack.v2

# 配置管理
go get -u gopkg.in/yaml.v2

# 密码加密
go get -u golang.org/x/crypto/bcrypt

# Swagger文档
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files



前端代码都在web文件里

服务端启动

go mod tidy

fresh


# 现阶段计划

CI/CD 流水线可视化平台
一个基于 Gin 框架的 CI/CD 流水线可视化平台后端，提供完整的流水线管理、构建、部署和监控功能。

功能特性
用户认证与授权：JWT 认证，角色权限控制
流水线管理：创建、配置、执行和监控流水线
构建管理：支持多种构建环境和构建模板
制品管理：管理构建产物和部署包
环境管理：管理不同的部署环境
发布管理：版本发布和回滚
API 文档：集成 Swagger 文档
日志系统：自定义彩色日志，支持日志分割和级别控制
数据库集成：GORM + MySQL，自动迁移
Redis 集成：支持缓存和会话管理
技术栈
Gin：高性能 Web 框架
GORM：ORM 库，支持自动迁移
JWT：用户认证
Zap：高性能日志库
Viper：配置管理
Swagger：API 文档
MySQL：数据库
Redis：缓存和会话管理
项目结构
``` gin_pipeline/ ├── api/ # API 接口 ├── config/ # 配置结构 ├── docs/ # 文档 ├── global/ # 全局变量 ├── initialize/ # 初始化 ├── logs/ # 日志文件 ├── middleware/ # 中间件 ├── model/ # 数据模型 ├── router/ # 路由 ├── service/ # 业务逻辑 ├── utils/ # 工具函数 ├── .gitignore # Git 忽略文件 ├── config.yaml # 配置文件 ├── go.mod # Go 模块文件 ├── go.sum # Go 模块依赖 ├── main.go # 入口文件 └── README.md # 说明文档 ```

快速开始
前置条件
Go 1.16+
MySQL 5.7+
Redis 5.0+
安装
克隆仓库
```bash git clone https://github.com/yourusername/gin_pipeline.git cd gin_pipeline ```

安装依赖
```bash go mod tidy ```

配置数据库
修改 config.yaml 文件中的数据库配置

运行项目
```bash go run main.go ```

或者使用热重载工具 fresh：

```bash go install github.com/pilu/fresh@latest fresh ```

访问 API 文档
``` http://localhost:8080/swagger/index.html ```

API 文档
项目集成了 Swagger 文档，运行项目后可以通过 /swagger/index.html 访问。

部署
Docker 部署
构建 Docker 镜像
```bash docker build -t gin_pipeline . ```

运行容器
```bash docker run -p 8080:8080 gin_pipeline ```

传统部署
编译项目
```bash go build -o gin_pipeline main.go ```

运行
```bash ./gin_pipeline ```

贡献
欢迎提交 Issue 和 Pull Request。

许可证
MIT ```

FROM golang:1.20-alpine AS builder

WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin_pipeline .

# 使用轻量级的 alpine 镜像
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /app/gin_pipeline .
COPY --from=builder /app/config.yaml .

# 创建必要的目录
RUN mkdir -p /app/logs /app/uploads

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./gin_pipeline"]