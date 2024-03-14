package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"visualon.com/go-server-monitor/config"
)

// StartHTTPServer Starts the webserver
func StartHTTPServer(port int) error {
	var err error

	go runCron()

	r := mux.NewRouter()

	r.PathPrefix("/monitor").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer w.(http.Flusher).Flush()
		log.Printf("Receiving stats from URL : %s", r.URL.String())
		switch r.Method {
		case http.MethodPost:
			PostMonitorHandler(r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})).Methods(http.MethodPost)

	r.PathPrefix("/logs").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer w.(http.Flusher).Flush()
		log.Printf("Receiving logs from URL : %s", r.URL.String())
		switch r.Method {
		case http.MethodPost:
			PostLogHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})).Methods(http.MethodPost)

	if (config.CONFIG.Secure.CertFilePath != "") && (config.CONFIG.Secure.KeyFilePath != "") {
		// Try HTTPS
		log.Printf("HTTPS server running on port %d", port)
		err = http.ListenAndServeTLS(":"+strconv.Itoa(port), config.CONFIG.Secure.CertFilePath, config.CONFIG.Secure.KeyFilePath, r)
	} else {
		// Try HTTP
		log.Printf("HTTP server running on port %d", port)
		err = http.ListenAndServe(":"+strconv.Itoa(port), r)
	}


	return err
}
