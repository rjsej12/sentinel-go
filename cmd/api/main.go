package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rjsej12/sentinel-go/internal/health"
	"github.com/rjsej12/sentinel-go/internal/metrics"
	"github.com/rjsej12/sentinel-go/internal/server"
	"github.com/rjsej12/sentinel-go/internal/worker"
)

func main() {
	metrics.Register()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queueSize := 100
	queue := worker.NewQueue(ctx, queueSize)

	workerPool := 5
	processor := worker.NewProcessor(ctx, queue, workerPool)
	processor.Start()

	router := server.NewRouter(queue, processor)
	handler := server.Logging(server.HTTPMetrics(router))
	httpServer := server.NewHTTPServer(":8080", handler)

	health.SetReady(true)
	httpServer.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	health.SetReady(false)
	processor.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped")
}
