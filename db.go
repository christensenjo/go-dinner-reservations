package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewDB(config DBConfig) *sql.DB {
	// continue from here
}
