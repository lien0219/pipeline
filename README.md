# 打造企业级标准化 CICD 平台

# 1.0.2 版本支持

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

1. **配置管理**：使用 YAML 配置文件，支持不同环境配置
2. **日志系统**：集成 zap 日志库，支持日志分割和级别控制
3. **数据库集成**：GORM + MySQL，包含连接池配置
4. **Redis 集成**：支持缓存和会话管理
5. **认证系统**：JWT 认证中间件
6. **API 文档**：集成 Swagger (swag)
7. **中间件**：日志、恢复、CORS 等
8. **路由管理**：分组和版本控制
9. **错误处理**：统一的错误响应格式

# 前端 (Vue 3)

1. **项目结构**：符合 Vue 3 最佳实践
2. **状态管理**：使用 Pinia
3. **路由**：Vue Router 配置，包含路由守卫
4. **API 层**：Axios 封装，拦截器配置
5. **认证**：完整的登录/注册/个人资料功能
6. **布局**：响应式布局组件
7. **页面**：首页、仪表盘、个人资料等

# Web 框架

go get -u github.com/gin-gonic/gin

# ORM 和数据库驱动

go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# Redis 客户端

go get -u github.com/go-redis/redis/v8

# JWT 认证

go get -u github.com/dgrijalva/jwt-go

# 日志库

go get -u go.uber.org/zap
go get -u gopkg.in/natefinch/lumberjack.v2

# 配置管理

go get -u gopkg.in/yaml.v2

# 热重载工具

go get -u github.com/pilu/fresh

# 跨域处理

go get -u github.com/gin-contrib/cors

# 错误处理

go get -u github.com/pkg/errors

# 响应处理

go get -u github.com/gin-contrib/response

# 验证器

go get -u github.com/go-playground/validator/v10

# 缓存

go get -u github.com/patrickmn/go-cache

# 工具

go get -u github.com/gin-contrib/sessions
go get -u github.com/gin-contrib/sessions/cookie
go get -u github.com/gin-contrib/sessions/redis
go get -u github.com/gin-contrib/sessions/memcache
go get -u github.com/gin-contrib/sessions/mongo
go get -u github.com/gin-contrib/sessions/postgres
go get -u github.com/gin-contrib/sessions/sqlite
go get -u github.com/gin-contrib/sessions/consul
go get -u github.com/gin-contrib/sessions/etcd
go get -u github.com/gin-contrib/sessions/mysql
go get -u github.com/gin-contrib/sessions/boltdb
go get -u github.com/gin-contrib/sessions/dynamodb
go get -u github.com/gin-contrib/sessions/gocache
go get -u github.com/gin-contrib/sessions/memcached
go get -u github.com/gin-contrib/sessions/mongodb
go get -u github.com/gin-contrib/sessions/negroni
go get -u github.com/gin-contrib/sessions/radix
go get -u github.com/gin-contrib/sessions/riak
go get github.com/shirou/gopsutil/v3

# 密码加密

go get -u golang.org/x/crypto/bcrypt

# Swagger 文档

go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

前端代码都在 web 文件里

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

```gin_pipeline/
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
`bash go mod tidy `

配置数据库
修改 config.yaml 文件中的数据库配置

运行项目
`bash go run main.go `

或者使用热重载工具 fresh：

`bash go install github.com/pilu/fresh@latest fresh `

访问 API 文档
`http://localhost:8080/swagger/index.html`

API 文档
项目集成了 Swagger 文档，运行项目后可以通过 /swagger/index.html 访问。

# resources/schemas/kubernetes/

`resources/schemas/kubernetes/` 目录下的 JSON 文件是 YAML 校验引擎功能的核心组件。这些文件是 JSON Schema 定义，用于验证 Kubernetes 资源的 YAML 配置是否符合规范。具体来说：

1. **`deployment.json`** - 这是 Kubernetes Deployment 资源的 JSON Schema 定义，用于验证 Deployment YAML 配置的正确性。它定义了 Deployment 资源必须包含的字段（如 apiVersion、kind、metadata、spec 等）以及各字段的数据类型、格式要求和约束条件。
2. **`service.json`** - 这是 Kubernetes Service 资源的 JSON Schema 定义，用于验证 Service YAML 配置的正确性。它定义了 Service 资源的必要字段和格式要求，如端口配置、选择器等。
3. **`general.json`** - 这是一个通用的 Kubernetes 资源 JSON Schema 定义，用于验证不在特定类型中的其他 Kubernetes 资源。它定义了所有 Kubernetes 资源共有的基本结构。

Schema 文件的作用是：

- **提供验证标准**：当用户提交 Kubernetes YAML 配置时，YAML 校验引擎会根据这些 Schema 定义验证配置的正确性
- **防止错误配置**：在部署前捕获配置错误，避免将错误的配置部署到生产环境
- **提供自动补全和提示**：可以基于这些 Schema 为用户提供编辑 YAML 时的字段提示和自动补全功能
  在 CI/CD 流程中，这些 Schema 文件被`yaml_validator.go`服务使用，以确保流水线中的 Kubernetes 配置在部署前是有效的，从而提高部署的成功率和系统的稳定性。

贡献
欢迎提交 Issue 和 Pull Request。

许可证
MIT ```

FROM golang:1.20-alpine AS builder
