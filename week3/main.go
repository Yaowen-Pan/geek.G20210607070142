package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	timeoutVal = 5 // set the maximum waiting time in seconds
)

// Listen to the interrupt signal of the system to stop the service
func listenAndShutdown(ctx context.Context, server *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-quit:
		fmt.Println("the service will stop")
	}
	return shutDownWithTimeout(server)
}

// provide a waiting time to avoid waiting indefinitely.
func shutDownWithTimeout(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutVal)*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	stopChannel := make(chan struct{})

	// define routes
	r := http.NewServeMux()
	r.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) { // responding to shutdown request
		stopChannel <- struct{}{}
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { // test A long link that takes one minute
		time.Sleep(time.Millisecond * 1000 * 60)
		fmt.Println("pong")
		w.Write([]byte("pong"))
	})

	// define server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// server start
	g.Go(func() error {
		fmt.Println("the service has been started")
		return server.ListenAndServe()
	})

	// server shutdown
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-stopChannel:
			fmt.Printf("the service will stop")
		}
		return shutDownWithTimeout(server)
	})

	// listen to the interrupt signal of the system
	g.Go(func() error {
		return listenAndShutdown(ctx, server)
	})

	// block
	if err := g.Wait(); err != nil {
		println(err)
	}
	println("the service has been stopped")
}
