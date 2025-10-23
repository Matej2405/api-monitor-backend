package database

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/mattn/go-sqlite3"
)
var DB *sql.DB

func InitDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "./api_monitor.db")
    if err != nil {
        return err
    }

    // Read and execute schema
    schema, err := os.ReadFile("internal/database/schema.sql")
    if err != nil {
        return err
    }

    _, err = DB.Exec(string(schema))
    if err != nil {
        return err
    }

    log.Println("Database initialized successfully")
    return nil
}

func CloseDB() {
    if DB != nil {
        DB.Close()
    }
}