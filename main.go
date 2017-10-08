package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"github.com/Sovianum/acquaintanceServer/config"
	"github.com/Sovianum/acquaintanceServer/server"
	"net/http"
	"github.com/gorilla/handlers"
)

const (
	confFile = "resources/config.json"
)

func main() {
	var conf = getConf()

	var db, dbErr = sql.Open(
		conf.DB.DriverName,
		conf.DB.GetAuthStr(),
	)
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	var env = server.NewEnv(db, conf)
	var router = server.GetRouter(env)
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, router))
}

func getConf() config.Conf {
	var file, confErr = os.Open(confFile)
	if confErr != nil {
		panic(confErr)
	}
	defer file.Close()

	var conf, err = config.ReadConf(file)
	if err != nil {
		panic(err)
	}

	return conf
}
