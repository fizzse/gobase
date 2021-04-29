package main

import (
	"context"
	"log"
	"time"

	"github.com/fizzse/gobase/internal/gobase/server"
)

func main() {
	app, cleanFunc, err := server.InitApp()
	if err != nil {
		log.Fatal("init app failed: ", err)
	}

	defer cleanFunc()
	if err := app.Run(context.Background()); err != nil {
		log.Println("recv error: ", err)
	}

	time.Sleep(2 * time.Second)
	log.Println("server graceful exit")
}
