package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "webPractice1/docs"
	"webPractice1/internal/repository"
	"webPractice1/internal/server"
	"webPractice1/internal/service"
	"webPractice1/internal/transport/handlers"
	"webPractice1/internal/transport/rabbitmq"
	"webPractice1/pkg/hasher"
	"webPractice1/pkg/logger"

	"github.com/spf13/viper"
)

// @title           rest with swagger and authorization with JWT tokens
// @version         1.0
// @description     project5
//
// @host      localhost:8080
// @BasePath  /
//
// JWT Authentication:
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
//
// RefreshToken authentication
// @securityDefinitions.apikey RefreshTokenAuth
// @in header
// @name RefreshToken
func init() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %v", err)
	}
}
func main() {
	logger := logger.GetLogger()
	hash := hasher.NewHashInit(viper.GetString("hashphrase"))
	repo := repository.NewRepository(repository.PostgresqlConnect(), logger)
	producer := rabbitmq.NewRabbitMQProducer(logger)
	service := service.NewService(repo, hash, logger, producer)
	handler := handlers.NewHandlerAssetsResponse(logger, service)
	srv := new(server.Server)
	go func() {
		err := srv.StartServer(handler.InitRoutes(), viper.GetString("server.port"))
		if err != nil {
			logger.Error(fmt.Sprintf("Server dont start: %s", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Shutdown error: %s", err))
	}
	log.Println("timeout of 1 seconds.")
	<-ctx.Done()

	log.Println("Server exiting")
}
