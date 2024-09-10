package main

import (
	"context"
	"log"
	"time"

	"github.com/AndresKenji/reverse-proxy/internal/server"
)

func main() {

	for {
		ctx, cancel := context.WithCancel(context.Background())
		srv := server.NewServer(ctx)
		cfgFile, err := srv.GetLatestConfig()
		if err != nil {
			log.Fatal(err.Error())
		}
		srv.SetServerMux(cfgFile)

		// Iniciar el servidor en una goroutine
		go func() {
			log.Println("API GateWay running on port:", srv.Port)
			if err := srv.StartServer(); err != nil {
				log.Panic(err)
			}
		}()
		// Timer para recarga de configuración
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		checkForRestart := make(chan bool)

		// Iniciar una goroutine que monitorea si la variable `restartServer` se vuelve true
		go func() {
			for range ticker.C {
				log.Println("Checking for configurations updates ...")
				latestCfg, err := srv.GetLatestConfig()
				if err != nil {
					log.Fatal(err.Error())
				}
				if latestCfg.CreatedAt.After(cfgFile.CreatedAt) {
					log.Println("Found new config")
					checkForRestart <- true
				}
			}
		}()

		// Esperar a que se envíe `true` para reiniciar el servidor o cancelar el contexto
		select {
		case <-checkForRestart:
			log.Println("Restarting the server...")
			cancel() // Cancelar el contexto actual para apagar el servidor
		case <-ctx.Done():
			log.Println("Context canceled, shutting down the server...")
			return // Salir si el contexto fue cancelado por otra razón
		}
	}
}
