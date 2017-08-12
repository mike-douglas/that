package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/mike-douglas/that/jsonpmiddleware"
	"github.com/mike-douglas/that/proxy"
)

func validate(next http.HandlerFunc) http.HandlerFunc {
	var authKey = os.Getenv("THAT_AUTH_KEY")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("key") == authKey {
			next(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()
	prox := prox.New(os.Getenv("THAT_COUCH_DB"))

	r.HandleFunc("/{that}/{thing}", validate(prox.Handle)).
		Methods(http.MethodDelete, http.MethodPost, http.MethodPut)
	r.HandleFunc("/{that}/{thing}", jsonpmiddleware.HandleJSONP(prox.Handle)).
		Methods(http.MethodGet)
	r.HandleFunc("/{that}", jsonpmiddleware.HandleJSONP(prox.Handle)).
		Methods(http.MethodGet)

	fmt.Printf("Listening on %s\n", ":8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
