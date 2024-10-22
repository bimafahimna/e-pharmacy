package main

import (
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/handler"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository/postgres"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/router"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/server"
)

func main() {
	config := config.InitConfig()

	logger.SetLogrusLogger(config.App)

	db, err := postgres.Init(config)
	if err != nil {
		logger.Log.Fatal("failed to connect db")
	}
	defer db.Close()

	opts := handler.Init(db, config)
	r := router.Init(opts, config)

	s := server.NewServer(config, r)
	s.Run()
}
