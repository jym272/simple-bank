package db

import (
	"database/sql"
	_ "github.com/lib/pq" // postgres driver
	"os"
	"simple_bank/utils"
	"testing"
)

var testQueries *Queries

type Connection struct {
	conn *sql.DB
}

var connection *Connection

func (c *Connection) DeleteAccountTable() {
	err := c.conn.QueryRow("TRUNCATE TABLE accounts CASCADE").Scan()
	if err != nil && err != sql.ErrNoRows {
		panic("cannot delete accounts table: " + err.Error())
	}
}

func (c *Connection) GetConnection() *sql.DB {
	return c.conn
}

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../")
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
	connection = &Connection{conn}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			panic("cannot close db: " + err.Error())
		}
	}(conn)

	testQueries = New(conn)
	os.Exit(m.Run())
}
