package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"messagio/internal/config"
	storage2 "messagio/internal/storage"
	"os"
)

func NewApp() {

	storage := storage2.NewConn()

	gin.SetMode(gin.ReleaseMode)

	handler := gin.Default()

}

func main() {
	config.In()
	os.Setenv("POSTGRES_URL", "sd")
	println(viper.GetString("POSTGRES.URL"))

}
