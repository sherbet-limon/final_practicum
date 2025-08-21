package main

import (
	"go1f/pkg/db"
	"go1f/pkg/server"
	"log"
)

const FormatDate string = "20060102"

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Printf("Ошибка инициализации БД: %v\n", err)
		return
	}
	defer func() {
        if db.DB != nil {
            db.DB.Close()
        }
    }()
	log.Println("БД успешно инициализирована")
	server.Run()
}
