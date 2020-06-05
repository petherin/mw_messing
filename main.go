package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {

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

func Handler() http.Handler {
	r := chi.NewRouter()

	r.With(
		mw1,
		mw2,
	).Handle("/endpoint", BaseHandler())

	return r
}

func BaseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("BaseHandler StatusMethodNotAllowed")

		//w.Write([]byte("{\"status\": \"running\"}"))
		w.WriteHeader(http.StatusMethodNotAllowed)
		// do nothing
	})
}
