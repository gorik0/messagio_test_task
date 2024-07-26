package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"messagio/internal/api"
	"messagio/internal/config"
	service2 "messagio/internal/service"
	storage2 "messagio/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerConfig struct {
	Port    string            // Порт, на котором запускается сервер
	Handler *gin.Engine       // Обработчик запросов Gin
	Handle  *api.Handler      // Обработчики API
	Service *service2.Service // Сервис для бизнес-логики
}

func (c ServerConfig) StartServer() {
	go func() {
		log.Printf("Starting server on %s\n", c.Port)
		if err := c.NewServer().ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (c ServerConfig) NewServer() *http.Server {

	return &http.Server{Addr: c.Port, Handler: c.Handler}
}

func (c ServerConfig) ShutdownServer() {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down server on %s\n", c.Port)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	if err := c.NewServer().Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	err := c.Service.CloseDB()
	if err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	println("GOOD bye!!!")
}

func NewApp() *ServerConfig {

	storage := storage2.NewConn()

	gin.SetMode(gin.ReleaseMode)

	gi := gin.Default()

	service := service2.NewService(storage)

	handler := api.NewHandler(gi, service)
	handler.SetupRoutes()

	service.StartConsumeMessages()

	return &ServerConfig{
		Port:    config.ServerAddress(),
		Handler: gi,
		Handle:  handler,
		Service: service,
	}
}

func main() {
	serverConfig := NewApp()

	serverConfig.StartServer()

	serverConfig.ShutdownServer()
}
