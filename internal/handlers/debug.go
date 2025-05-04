package handlers

import (
	"fmt"
	"net/http"
)

func DiagnosticsPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Debug Information</h1>"))

	// Show all headers
	w.Write([]byte("<h2>Headers</h2><pre>"))
	for name, values := range r.Header {
		for _, value := range values {
			w.Write([]byte(fmt.Sprintf("%s: %s\n", name, value)))
		}
	}
	w.Write([]byte("</pre>"))

	// Show cookies
	w.Write([]byte("<h2>Cookies</h2><pre>"))
	for _, cookie := range r.Cookies() {
		w.Write([]byte(fmt.Sprintf("%s: %s\n", cookie.Name, cookie.Value)))
	}
	w.Write([]byte("</pre>"))
}
