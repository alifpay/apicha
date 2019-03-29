package main

import (
	"alif/apitemp/config"
	"alif/apitemp/db"
	"alif/apitemp/handlers"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		configPath = flag.String("config", "./config.json", "path of the config file")
		err        error
	)
	flag.Parse()
	//parse configs
	cfg, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatalln("config error", err)
		return
	}

	// connecting to datadase
	db.Connect(cfg.Database.Addr, cfg.Database.User, cfg.Database.Pass, cfg.Database.DBName)
	defer db.Close()
	log.Println("db is connected")

	// declaring our routes
	http.HandleFunc("/ping", handlers.Ping)

	ctx, cancelFun := context.WithCancel(context.Background())

	srv := &http.Server{
		Addr:         cfg.Service.Addr,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	go func() {
		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
		s := <-sigint
		log.Println("server received signal ", s)
		defer cancelFun()
		err = srv.Shutdown(ctx)
		if err != nil {
			log.Println("server: couldn't shutdown because of " + err.Error())
		}
	}()
	log.Println("http server is runnning", cfg.Service.Addr)
	log.Fatal(srv.ListenAndServe())
}
