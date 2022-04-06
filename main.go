package main

import (
	"fmt"
	"net/http"
)

// testHandler is a struct with no fields declared only to associate the ServeHTTP method
type testHandler struct{}

func (h testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello client from Handle")
}

func testMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from testMW\n")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. http.Handle -- we use the ServeHTTP method from testHandler empty struct as the request handler
	http.Handle("/test-handle", testHandler{})

	// 2. http.Handle -- we use the HandlerFunc type to avoid declaring a type for each handler
	testFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello client from HandlerFunc type")
	}

	http.Handle("/test-func", http.HandlerFunc(testFunc))

	// 3. http.HandleFunc -- we use an anonymous function as the request handler
	http.HandleFunc("/test-handlefunc", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello client from HandleFunc")
	})

	// 4. handling GET requests
	testGET := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// iterate over form values
			for key, val := range r.Form {
				fmt.Printf("key ==> %v val ==> %v\n", key, val)
			}

			// ask for a specific key
			testValue := r.FormValue("test")
			fmt.Printf("testValue ==> %v\n", testValue)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}

	http.HandleFunc("/test-get", testGET)

	// 5. How to implement middleware functions
	testMWFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from testMWFunc")
	})

	http.Handle("/test-mw", testMW(testMWFunc))

	// 6. serving assets
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/", fs)

	// initialize the HTTP Server
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
