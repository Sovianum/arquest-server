package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Configuration struct {
	Port       int
	DriverName string
	DBUser     string
	DBPassword string
	DBName     string
	LogFile    string
}

func main() {
	var file, confErr = os.Open("conf.json")
	if confErr != nil {
		panic(confErr)
	}
	defer file.Close()
	var conf = Configuration{}

	var parseErr = json.NewDecoder(file).Decode(&conf)
	if parseErr != nil {
		panic(parseErr)
	}

	var db, err = sql.Open(
		conf.DriverName,
		getDBStr(conf.DBUser, conf.DBPassword, conf.DBName),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//var logFileWriter, logErr = os.Create(conf.LogFile)
	//if logErr != nil {
	//	panic(logErr)
	//}
	//defer logFileWriter.Close()
}

func getDBStr(user string, pass string, name string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, name)
}
