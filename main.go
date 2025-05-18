package main

import (
	"context"
	"flag"
	"fmt"
	"gin_pipeline/config"
	"gin_pipeline/docs"
	"gin_pipeline/global"
	"gin_pipeline/initialize"
	"gin_pipeline/middleware"
	"gin_pipeline/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Enterprise API
// @version         1.0
// @description     Enterprise API Server
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	if err := config.Init(*configPath); err != nil {
		panic(fmt.Sprintf("Failed to initialize configuration: %v", err))
	}

	if config.Get().Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if config.Get().Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if err := initialize.Logger(); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	if err := initialize.MySQL(); err != nil {
		global.Logger.Error("Failed to initialize MySQL: " + err.Error())
		panic(err)
	}

	if err := initialize.Redis(); err != nil {
		global.Logger.Error("Failed to initialize Redis: " + err.Error())
		panic(err)
	}

	if err := initialize.MySQL(); err != nil {
		global.Logger.Error("Failed to initialize MySQL: " + err.Error())
		panic(err)
	}

	// 执行数据库迁移
	if err := initialize.Migrate(); err != nil {
		global.Logger.Error("Failed to migrate database: " + err.Error())
		panic(err)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"

	// Create Gin engine
	app := gin.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())
	app.Use(middleware.Cors())

	router.RegisterRoutes(app)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Get().Server.Port),
		Handler: app,
	}

	go func() {
		global.Logger.Info(fmt.Sprintf("Starting server on port %d", config.Get().Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error("Failed to start server: " + err.Error())
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		global.Logger.Error("Server forced to shutdown: " + err.Error())
	}

	global.Logger.Info("Server exiting")
}
