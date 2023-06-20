package main

import (
	"database/sql"
	_ "github.com/lib/pq" // postgres driver
	"simple_bank/api"
	db "simple_bank/sqlc"
	"simple_bank/utils"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		panic("cannot load config: " + err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	err = conn.Ping()
	if err != nil {
		panic("cannot ping db: " + err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		panic("cannot start server: " + err.Error())
	}

}
