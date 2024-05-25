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

	// Inisialisasi koneksi PostgreSQL
	postgresConn, errPostgres := config.ConnectPostgreSQL()
	if errPostgres != nil {
		log.Panic("error PostgreSQL connection: ", errPostgres)
	}
	defer func() {
		if sqlDB, err := postgresConn.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Inisialisasi koneksi Redis
	redisClient := config.NewRedisClient()

	// Inisialisasi router dengan koneksi PostgreSQL dan Redis
	router := app.InitRouter(postgresConn, redisClient)
	log.Println("routes initialized")

	// Mendapatkan port dari konfigurasi
	port := viper.GetString("PORT")

	// Inisialisasi server HTTP
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Println("Server initialized, listening at port:", port)

	// Mulai server HTTP dalam goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err.Error())
		}
	}()

	// Tunggu sinyal shutdown dari OS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	// Set timeout context untuk shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
