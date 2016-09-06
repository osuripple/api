package main

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"git.zxq.co/ripple/rippleapi/app"
	"git.zxq.co/ripple/rippleapi/beatmapget"
	"git.zxq.co/ripple/rippleapi/common"
	"git.zxq.co/ripple/schiavolib"
	// Golint pls dont break balls
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/thehowl/go-osuapi.v1"
)

// Version is the git hash of the application. Do not edit. This is
// automatically set using -ldflags during build time.
var Version string

func init() {
	log.SetFlags(log.Ltime)
	log.SetPrefix(fmt.Sprintf("%d|", syscall.Getpid()))
	common.Version = Version
}

var db *sqlx.DB

func main() {
	fmt.Print("Ripple API")
	if Version != "" {
		fmt.Print("; git commit hash: ", Version)
	}
	fmt.Println()

	conf, halt := common.Load()
	if halt {
		return
	}

	schiavo.Prefix = "Ripple API"

	if !strings.Contains(conf.DSN, "parseTime=true") {
		c := "?"
		if strings.Contains(conf.DSN, "?") {
			c = "&"
		}
		conf.DSN += c + "parseTime=true"
	}

	var err error
	db, err = sqlx.Open(conf.DatabaseType, conf.DSN)
	if err != nil {
		schiavo.Bunker.Send(err.Error())
		log.Fatalln(err)
	}

	beatmapget.Client = osuapi.NewClient(conf.OsuAPIKey)
	beatmapget.DB = db

	engine := app.Start(conf, db)

	startuato(engine)
}
