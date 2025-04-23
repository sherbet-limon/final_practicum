package main

import (
	"go1f/pkg/server"
	"go1f/pkg/db"
	// "go1f/pkg/api"
	"log"
)

const FormatDate string = "20060102"

func main() {
	if err:= db.Init("scheduler.db"); err != nil {
		log.Printf("Ошибка инициализации БД: %v\n", err)
	}
	log.Println("БД успешно инициализирована")
	server.Run()
	
}
