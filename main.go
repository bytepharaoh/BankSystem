package main

import (
	"database/sql"
	"log"

	"github.com/bytepharoh/simplebank/api"
	db "github.com/bytepharoh/simplebank/db/sqlc"
	"github.com/bytepharoh/simplebank/util"
	_ "github.com/lib/pq"
)


func main() {
	config , err :=util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can not establish a connection:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAdress)
	if err != nil {
		log.Fatal("can not start the server!")
	}
}
