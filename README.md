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