package main

import (
	"go1f/pkg/server"
	"go1f/pkg/db"
	"log"
)
func main() {
	if err:= db.Init("scheduler.db"); err != nil {
		log.Printf("Ошибка инициализации БД: %v\n", err)
	}
	log.Println("БД успешно инициализирована")
	server.Run()
	
}
