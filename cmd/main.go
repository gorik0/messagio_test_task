package main

import (
	"github.com/spf13/viper"
	"messagio/internal/config"
)

func main() {
	println(viper.GetString("PostgresUrl"))
	config.In()
}
