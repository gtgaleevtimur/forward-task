// Package app - аккумулирует все пакеты сервиса для запуска приложения.
package app

import (
	"context"
	"fmt"
	"forward/internal/config"
	"forward/internal/handler"
	"forward/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run - запускает сервис.
func Run() {
	conf := config.NewConfig()
	db := repository.NewRepository()

	switch conf.Mode {
	case "update-registry":
		err := db.UpdateRegistry()
		if err != nil {
			log.Fatal(err)
		}
	case "serve":
		if conf.ServerAddress == "" {
			fmt.Println("Не задан адрес сервера.")
			return
		}
		startServe(conf, db)
	default:
		fmt.Println("Режим работы не задан.")
	}
}

// startServe - запускает сервер с Grace-Ful Shutdown.
func startServe(conf *config.Config, db *repository.Repository) {
	server := &http.Server{
		Addr:    conf.ServerAddress,
		Handler: handler.NewRouter(db),
	}
	go func() {
		log.Println("starting server at:", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	// Канал Grace-ful Shutdown.
	sigs := make(chan os.Signal)
	signal.Notify(sigs,
		syscall.SIGINT,
		os.Interrupt)
	// После получения сигнала закрываем приложение.
	<-sigs
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("server shutdown error")
	}
	log.Println("shutting down")
	os.Exit(0)
}
