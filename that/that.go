package main

import (
	"fmt"

	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/go-uuid"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/mike-douglas/that/jsonpmiddleware"
	"github.com/mike-douglas/that/proxy"
)

func generateUUID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUUID, _ = uuid.GenerateUUID()

		r.URL, _ = url.Parse(fmt.Sprintf("%s/%s", r.URL.String(), newUUID))

		next(w, r)
	}
}

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
	fmt.Printf("Key: %s\nDB: %s\n\n", os.Getenv("THAT_AUTH_KEY"), os.Getenv("THAT_COUCH_DB"))

	r := mux.NewRouter()
	prox := prox.New(os.Getenv("THAT_COUCH_DB"))

	r.HandleFunc("/{that}/{thing}", validate(prox.Handle)).
		Methods(http.MethodDelete, http.MethodPost, http.MethodPut)
	r.HandleFunc("/{that}/{thing}", jsonpmiddleware.HandleJSONP(prox.Handle)).
		Methods(http.MethodGet)
	r.HandleFunc("/{that}", validate(prox.Handle)).
		Methods(http.MethodDelete, http.MethodPut)
	r.HandleFunc("/{that}", validate(generateUUID(prox.Handle))).
		Methods(http.MethodPost)
	r.HandleFunc("/{that}", jsonpmiddleware.HandleJSONP(prox.Handle)).
		Methods(http.MethodGet)

	fmt.Printf("Listening on %s\n", ":8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
