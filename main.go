package main

import (
	"database/sql"
	_ "github.com/lib/pq" // postgres driver
	"simple_bank/api"
	db "simple_bank/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://postgres:postgres@localhost:8080/simple_bank?sslmode=disable"
	serverAddress = ":8081"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	err = conn.Ping()
	if err != nil {
		panic("cannot ping db: " + err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		panic("cannot start server: " + err.Error())
	}

}
