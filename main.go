package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/m/handlers"
)

func main() {
	l := log.New(os.Stdout, "whoami ", log.LstdFlags)

	wh := handlers.NewWhoami(l)

	sm := http.NewServeMux()
	sm.Handle("/api/whoami", wh)

	s := &http.Server{
		Addr:         ":" + getPort(),
		Handler:      sm,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, c := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
	c()
}

func getPort() string {
	p := os.Getenv("PORT")

	if os.Getenv("PORT") == "" {
		p = "8080"
	}

	return p
}
