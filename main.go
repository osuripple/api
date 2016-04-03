package main

import (
	"database/sql"
	"log"

	// Golint pls dont break balls
	_ "github.com/go-sql-driver/mysql"
	"github.com/osuripple/api/app"
	"github.com/osuripple/api/common"
)

func main() {
	conf, halt := common.Load()
	if halt {
		return
	}
	db, err := sql.Open(conf.DatabaseType, conf.DSN)
	if err != nil {
		log.Fatal(err)
	}
	app.Start(conf, db)
}
