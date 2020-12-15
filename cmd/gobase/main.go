package main

import (
	"log"

	"github.com/fizzse/gobase/internal/gobase/server"
)

func main() {
	app, clanFunc, err := server.InitApp()
	if err != nil {
		log.Fatal("init app failed: ", err)
	}

	defer clanFunc()
	app.RestServer.ListenAndServe()
}
