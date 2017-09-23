package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"github.com/Sovianum/acquaintanceServer/config"
)

const (
	dbConfFile = "resources/conf/db_conf.json"
	authConfFile = "resources/conf/auth_conf.json"
)

func main() {
	var dbConf = getDBConf()
	//var authConf = getAuthConf()

	var db, dbErr = sql.Open(
		dbConf.DriverName,
		dbConf.GetAuthStr(),
	)
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	//var r = mux.NewRouter()
	//http.ListenAndServe(":3000", server.LoggingHandler(os.Stdout, r))

	//var logFileWriter, logErr = os.Create(dbConf.LogFile)
	//if logErr != nil {
	//	panic(logErr)
	//}
	//defer logFileWriter.Close()
}

func getAuthConf() config.AuthConfig {
	var file, confErr = os.Open(authConfFile)
	if confErr != nil {
		panic(confErr)
	}
	defer file.Close()

	var conf, err = config.ReadAuthConf(file)
	if err != nil {
		panic(err)
	}

	return conf
}

func getDBConf() config.DBConfig {
	var file, confErr = os.Open(dbConfFile)
	if confErr != nil {
		panic(confErr)
	}
	defer file.Close()

	var conf, err = config.ReadDBConfig(file)
	if err != nil {
		panic(err)
	}

	return conf
}
