package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// testHandler is a struct with no fields declared only to associate the ServeHTTP method
type testHandler struct{}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello client from Handle")
}

func main() {
	// 1. simple handler definition.
	r := mux.NewRouter()
	r.Handle("/test-handle", testHandler{})

	// 2. simple handle func.
	r.HandleFunc("/handle-func", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello handle func")
	})

	// 3. path with variables
	r.HandleFunc("/handle-vars/{vars}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Fprintf(w, "hello vars %s", vars["vars"])
	})

	// 4. routing with matchers.
	r.HandleFunc("/handle-match", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello matcher")
	}).Methods("GET").
		Name("matcher")

	// 5. subrouters.
	s := r.PathPrefix("/sub").Subrouter()
	// 5.1 GET /sub
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello sub")
	})

	// 5.2 GET /sub/one
	s.HandleFunc("/one", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello sub one")
	}).Methods("GET")

	// 6 Middleware.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
