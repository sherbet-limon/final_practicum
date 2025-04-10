package db

import ( 
	"database/sql"
	_ "modernc.org/sqlite"
)
	
const schema = `
CREATE TABL scheduler(
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT "",
title VARCHAR(256) NOT NULL DEFAULT "",
comment TEXT NOT NULL DEFAULT "",
repeat VARCHAR(128) NOT NULL DEFAULT "",
);`
const schemaIdx = `
CREATE INDEX date_id ON scheduler (date);`

