# 打造企业级标准化CICD平台

# 1.0.2版本支持

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
  日志系统：自定义日志，支持日志分割和级别控制
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

快速开始
前置条件
Go 1.16+
MySQL 5.7+
Redis 5.0+
安装
克隆仓库
```bash git clone https://github.com/lien0219/pipeline.git
   cd gin_pipeline
```
项目结构
``` gin_pipeline/ 
├── api/ # API 接口
├── config/ # 配置结构
├── docs/ # 文档
├── global/ # 全局变量
├── initialize/ # 初始化
├── logs/ # 日志文件
├── middleware/ # 中间件
├── model/ # 数据模型
├── router/ # 路由
├── service/ # 业务逻辑
├── utils/ # 工具函数
├── .gitignore # Git 忽略文件
├── config.yaml # 配置文件
├── go.mod # Go 模块文件
  ├── go.sum # Go 模块依赖
├── main.go # 入口文件
└── README.md # 说明文档
```
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

# resources/schemas/kubernetes/
`resources/schemas/kubernetes/` 目录下的JSON文件是YAML校验引擎功能的核心组件。这些文件是JSON Schema定义，用于验证Kubernetes资源的YAML配置是否符合规范。具体来说：

1. **`deployment.json`** - 这是Kubernetes Deployment资源的JSON Schema定义，用于验证Deployment YAML配置的正确性。它定义了Deployment资源必须包含的字段（如apiVersion、kind、metadata、spec等）以及各字段的数据类型、格式要求和约束条件。
2. **`service.json`** - 这是Kubernetes Service资源的JSON Schema定义，用于验证Service YAML配置的正确性。它定义了Service资源的必要字段和格式要求，如端口配置、选择器等。
3. **`general.json`** - 这是一个通用的Kubernetes资源JSON Schema定义，用于验证不在特定类型中的其他Kubernetes资源。它定义了所有Kubernetes资源共有的基本结构。

Schema文件的作用是：
- **提供验证标准**：当用户提交Kubernetes YAML配置时，YAML校验引擎会根据这些Schema定义验证配置的正确性
- **防止错误配置**：在部署前捕获配置错误，避免将错误的配置部署到生产环境
- **提供自动补全和提示**：可以基于这些Schema为用户提供编辑YAML时的字段提示和自动补全功能
在CI/CD流程中，这些Schema文件被`yaml_validator.go`服务使用，以确保流水线中的Kubernetes配置在部署前是有效的，从而提高部署的成功率和系统的稳定性。

贡献
欢迎提交 Issue 和 Pull Request。

许可证
MIT ```

FROM golang:1.20-alpine AS builder
