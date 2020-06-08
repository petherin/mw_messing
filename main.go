package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	handlers:=NewHandlers()

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlers,
		// These timeouts configure how long to keep connections open while reading, writing and idling
		// The IdleTimeout should be higher than the load balancer for keepalive to work correctly.
		ReadTimeout:       time.Second * 20,
		WriteTimeout:      time.Second * 20,
		ReadHeaderTimeout: time.Second * 20,
		IdleTimeout:       time.Second * 60,
	}

	errCh := make(chan error)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := http.ListenAndServe(server.Addr, server.Handler); err != nil {
			log.Println(fmt.Sprintf("server stopped due to error: %s"),err.Error())
			errCh <- err
		}
	}()

	// Let us hang here until we either receive a signal from the OS, or the app decides to shut itself down.
	select {
	case <-quit:
		log.Println("server shutting down, signal received from os")
	case err:= <-errCh:
		log.Println("server shutting down, error received on application error channel: %s", err.Error())
	}
}


func mw1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("mw1 StatusBadRequest")
		//w.Write([]byte("{\"mw1\": \"pow\"}"))
		//w.WriteHeader(http.StatusBadRequest)
		next.ServeHTTP(w, r)
	})
}

func mw2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("mw2 StatusInternalServerError")
		//w.Write([]byte("{\"mw2\": \"pow\"}"))
		//w.WriteHeader(http.StatusInternalServerError)
		next.ServeHTTP(w, r)
	})
}

func NewHandlers() http.Handler {
	r := chi.NewRouter()

	r.With(
		mw1,
		mw2,
	).Handle("/", BaseHandler())

	return r
}

func BaseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("BaseHandler StatusMethodNotAllowed")

		w.Write([]byte("{\"status\": \"running\"}"))
		//w.WriteHeader(http.StatusMethodNotAllowed)
		// do nothing
	})
}
