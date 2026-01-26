package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rjsej12/sentinel-go/internal/health"
	"github.com/rjsej12/sentinel-go/internal/metrics"
	"github.com/rjsej12/sentinel-go/internal/server"
)

func main() {
	metrics.Register()

	router := server.NewRouter()

	handler := server.Logging(metrics.HTTPMetrics(router))

	httpServer := server.NewHTTPServer(":8080", handler)

	health.SetReady(true)
	httpServer.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	health.SetReady(false)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
}
