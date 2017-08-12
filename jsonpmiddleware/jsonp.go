package jsonpmiddleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

/*HandleJSONP is middleware to convert a response to JSONP if a callback is present*/
func HandleJSONP(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if callback := r.FormValue("callback"); callback != "" {
			rr := httptest.NewRecorder()

			next(rr, r)

			if status := rr.Code; status != http.StatusOK {
				http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			fmt.Fprintf(w, "%s(%x)", callback, rr.Body.String())
		} else {
			next(w, r)
		}
	}
}
