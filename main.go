package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workspace-server/app/controller"
	"workspace-server/app/helper"
	"workspace-server/app/middleware"
	other_repository "workspace-server/app/repository/other"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/app/service"
	"workspace-server/config"
	server_grpc "workspace-server/grpc"
	"workspace-server/package/database"
	logger "workspace-server/package/log"
	_validator "workspace-server/package/validator"
	"workspace-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jutimi/grpc-service/workspace"
	"google.golang.org/grpc"
)

func main() {
	conf := config.GetConfiguration()

	// Register repositories
	postgresDB := database.GetPostgres()
	// mysqlRepo := mysql_repository.RegisterMysqlRepositories(db)
	postgresRepo := postgres_repository.RegisterPostgresRepositories(postgresDB)
	repo := other_repository.RegisterOtherRepositories()

	// Register Others
	helpers := helper.RegisterHelpers(postgresRepo, repo)
	services := service.RegisterServices(helpers, postgresRepo, repo)
	middleware := middleware.RegisterMiddleware()

	// Run GRPC Server
	go startGRPCServer(conf, postgresRepo, helpers)

	// Run gin server
	gin.SetMode(conf.Server.Mode)
	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logger.GetLogger().Writer()))

	// Register validator
	v := binding.Validator.Engine().(*validator.Validate)
	v.SetTagName("validate")
	_validator.RegisterCustomValidators(v)

	// Register controllers
	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})
	controller.RegisterControllers(router, services, middleware)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: router,
	}

	if !gin.IsDebugging() {
		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		// catching ctx.Done(). timeout of 5 seconds.
		select {
		case <-ctx.Done():
			log.Println("timeout of 5 seconds.")
		}
		log.Println("Server exiting")
	} else {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}
}

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.InitPostgres()
	logger.Init()
}

func startGRPCServer(
	conf *config.Configuration,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
	helpers helper.HelperCollections,
) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPC.OAuthPort))
	if err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	// Register server
	workspace.RegisterUserWorkspaceRouteServer(grpcServer, server_grpc.NewGRPCServer(postgresRepo, helpers))
	workspace.RegisterWorkspaceRouteServer(grpcServer, server_grpc.NewGRPCServer(postgresRepo, helpers))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error Init GRPC: %s", err.Error())
	}

	log.Println("Init GRPC Success!")
}
