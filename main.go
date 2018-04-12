package main

import (
	"database/sql"
	"fmt"
	"github.com/Sovianum/arquest-server/config"
	"github.com/Sovianum/arquest-server/mylog"
	"github.com/Sovianum/arquest-server/routes"
	"github.com/Sovianum/arquest-server/server"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"strconv"
)

const (
	confFile = "resources/config.json"
)

func main() {
	logger := mylog.NewLogger(os.Stdout)

	conf, err := getConf()
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	db, err := connectDB(conf, logger)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	env := server.NewEnv(db, conf, logger)
	router := routes.GetEngine(env)

	portLine := fmt.Sprintf(":%d", getServerPort(conf, logger))
	http.ListenAndServe(portLine, handlers.LoggingHandler(os.Stdout, router))
}

func getServerPort(conf *config.Conf, logger *mylog.Logger) int {
	portStr := os.Getenv(conf.PortEnvVar)

	if port, err := strconv.Atoi(portStr); err == nil {
		logger.Info("Used system port")
		return port
	}
	logger.Info("Used default port")
	return conf.DefaultPort
}

func connectDB(conf *config.Conf, logger *mylog.Logger) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if db, err = sql.Open(conf.DB.DriverName, conf.DB.GetEnvAuthString()); err != nil {
		return nil, err
	}
	if err = db.Ping(); err == nil {
		logger.Info("Authorized via env")
		return db, nil
	}
	if db, err = sql.Open(conf.DB.DriverName, conf.DB.GetAuthStr()); err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	logger.Info("Authorized via conf")
	return db, err
}

func getConf() (*config.Conf, error) {
	file, confErr := os.Open(confFile)
	if confErr != nil {
		return nil, confErr
	}
	defer file.Close()

	return config.ReadConf(file)
}
