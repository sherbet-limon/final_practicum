package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler(
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT "",
title VARCHAR(256) NOT NULL DEFAULT "",
comment TEXT NOT NULL DEFAULT "",
repeat VARCHAR(128) NOT NULL DEFAULT ""
);`
const schemaIdx = `
CREATE INDEX IF NOT EXISTS date_id ON scheduler (date);`

var DB *sql.DB

// проверяем наличие бд, если нет, то CREATE TABLE INDEX
func Init(dbFile string) error {
	var install bool
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть БД: %w", err)
	}
	//сохраняем открытое соединение
	DB = db
	if install {
		if _, err := db.Exec(schema); err != nil {
			log.Println("не удалось создать таблицу: %w", err)
		}
		if _, err := db.Exec(schemaIdx); err != nil {
			log.Println("не удалось создать индекс к date: %w", err)
		}
	}
	return nil
}
