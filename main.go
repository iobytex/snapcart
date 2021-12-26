package main

import (
	"log"
	"snapcart/config"
	"snapcart/model"
	"snapcart/pkg/postgres"
	"snapcart/server"
)

func main() {

	getConfig, err := config.GetConfig()
	if err != nil {
		return
	}
	initDBConnection, err := postgres.InitDBConnection(getConfig)
	if err != nil {
		panic(err)
	}

	migrateErr := initDBConnection.AutoMigrate(&model.Message{})
	if migrateErr != nil {
		log.Print(migrateErr)
	}


	s := server.NewServer(initDBConnection)
	sErr := s.Run()
	if sErr != nil {
		log.Print(sErr)
	}
}
