package main

import (
	"context"
	"log"
	"github.com/AndresKenji/reverse-proxy/internal/server"
)

func main() {

	for {
		checkForRestart := make(chan bool)
		ctx, cancel := context.WithCancel(context.Background())
		srv := server.NewServer(ctx, checkForRestart)
		cfgFile, err := srv.GetLatestConfig()
		if err != nil {
			log.Fatal(err.Error())
		}
		if cfgFile == nil {
			cfgFile = srv.SetDefaultConfig()
		}
		srv.SetServerMux(cfgFile)

		// Iniciar el servidor en una goroutine
		go func() {
			log.Println("API GateWay running on port:", srv.Port)
			if err := srv.StartServer(); err != nil {
				log.Panic(err)
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
