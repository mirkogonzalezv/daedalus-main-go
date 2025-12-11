package main

import (
	"daedalus-engine-go/internal/app"
	"log"
)

func main() {
	application := app.NewApp()

	if err := application.Initialize(); err != nil {
		log.Fatal("Error inicializando aplicación:", err)
	}
	defer application.Shutdown()

	if err := application.Run(); err != nil {
		log.Fatal("Error ejecutando aplicación:", err)
	}
}
