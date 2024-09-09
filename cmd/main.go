package main

import (
	"context"
	"log"
	"time"

	"github.com/AndresKenji/reverse-proxy/internal/server"
)

var restartServer bool

func main() {
	
	for {
		// Crear un contexto con cancelación manual
		ctx, cancel := context.WithCancel(context.Background())

		// Crear el servidor
		srv := server.NewServer(ctx)

		cfgFile, err := srv.GetLatestConfig()
		if err != nil {
			log.Fatal(err.Error())
		}

		// Configurar las rutas del servidor
		srv.SetServerMux(cfgFile)

		// Iniciar el servidor en una goroutine
		go func() {
			log.Println("API GateWay running on port:", srv.Port)
			if err := srv.StartServer(); err != nil {
				log.Panic(err)
			}
		}()

		// Monitorear la condición cada 10 segundos
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		checkForRestart := make(chan bool)

		// Iniciar una goroutine que monitorea si la variable `restartServer` se vuelve true
		go func() {
			for range ticker.C {
				restartServer = true
				if restartServer {
					checkForRestart <- true
					restartServer = false // Reiniciar la variable después de la verificación
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
