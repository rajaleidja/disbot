package database

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func Connect(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    log.Println("Connected to the database")
    return db, nil
}