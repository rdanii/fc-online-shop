package main

import (
	"flag"
	"log"
	"online-shop/config"
)

func main() {
	envF := flag.String("env", "local", "define environment")
	flag.Parse()

	env := *envF
	config.InitConfig(env)

	postgresConn, errPostgres := config.ConnectPostgreSQL()
	if errPostgres != nil {
		log.Panic("error PostgreSQL connection: ", errPostgres)
	}
	defer func() {
		if sqlDB, err := postgresConn.DB(); err == nil {
			sqlDB.Close()
		}
	}()
}
