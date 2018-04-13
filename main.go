package main

import (
	"database/sql"
	"fmt"
	"github.com/Sovianum/arquest-server/config"
	"github.com/Sovianum/arquest-server/mylog"
	"github.com/Sovianum/arquest-server/routes"
	"github.com/Sovianum/arquest-server/server"
	"github.com/Sovianum/arquest-server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	confFile   = "resources/ard.conf.json"
	defaultLog = "/var/log/ard.log"
)

func main() {
	flags := utils.NewFlags(confFile)
	flags.Parse()

	fmt.Printf("ARD started.\nConfig from %s\n", flags.Config)

	conf, err := getConf(flags)
	if err != nil {
		panic(err)
	}

	if conf.Log == "" {
		conf.Log = defaultLog
	}

	fmt.Printf("Logging to %s\n", conf.Log)

	logger, f, err := getLogger(conf)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gin.DefaultWriter = io.MultiWriter(f)

	db, err := connectDB(conf, logger)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	env := server.NewEnv(db, conf, logger)
	router := routes.GetEngine(env)

	portLine := fmt.Sprintf(":%d", getServerPort(conf, logger))
	if err := http.ListenAndServe(portLine, handlers.LoggingHandler(os.Stdout, router)); err != nil {
		panic(err)
	}
}

func getLogger(conf *config.Conf) (*mylog.Logger, *os.File, error) {
	_, err := os.Stat(conf.Log)

	var f *os.File
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(conf.Log)
			if err != nil {
				return nil, nil, err
			}
		}
	} else {
		var innerErr error
		f, innerErr = os.OpenFile(conf.Log, os.O_APPEND|os.O_WRONLY, 0600)
		if innerErr != nil {
			return nil, nil, innerErr
		}
	}
	return mylog.NewLogger(f), f, nil
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
	authStr := conf.DB.GetAuthStr()
	fmt.Printf("connecting to db via %s\n", authStr)
	if db, err = sql.Open(conf.DB.DriverName, authStr); err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	logger.Info("Authorized via conf")
	return db, err
}

func getConf(flags *utils.Flags) (*config.Conf, error) {
	file, confErr := os.Open(flags.Config)
	if confErr != nil {
		return nil, confErr
	}
	defer file.Close()

	return config.ReadConf(file)
}
