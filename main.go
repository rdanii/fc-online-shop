package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"online-shop/app"
	"online-shop/config"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
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

	router := app.InitRouter(postgresConn)
	log.Println("routes initialized")

	port := viper.GetString("PORT")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Println("Server initialized, listening at port:", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
