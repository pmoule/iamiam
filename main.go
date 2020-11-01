package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/pmoule/iamiam/iamiam"
)

const (
	configFileName string = "config/iam_config.json"
)

type configuration struct {
	Hostname       string
	Port           int
	ValidUserInfos []*iamiam.UserInfo
}

func main() {
	config, err := readConfiguration()

	if err != nil {
		log.Fatal(err)
	}

	router := InitRouter(config.ValidUserInfos)
	address := fmt.Sprintf("%s:%d", config.Hostname, config.Port)
	run(router, address)
}

func run(router *mux.Router, address string) {
	var wait time.Duration = time.Second * 1
	srv := &http.Server{
		Addr:         address,
		WriteTimeout: wait,
		ReadTimeout:  wait,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		log.Printf("running on address: %s", address)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// prepare graceful shutdown
	signal.Notify(c, os.Interrupt)

	// Block until signal.
	<-c
	log.Printf("prepare for shutdown (timeout: %d seconds)\n", int(wait.Seconds()))
	// wait for timeout
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// wait for timeout
	srv.Shutdown(ctx)
	log.Println("shutdown")
	os.Exit(0)
}

func readConfiguration() (*configuration, error) {
	log.Println("reading configuration")
	file, err := os.Open(configFileName)

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	config := configuration{}
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
