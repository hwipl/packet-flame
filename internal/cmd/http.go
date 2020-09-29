package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
)

// handleRequest returns the flamegraph to http clients
func handleRequest(w http.ResponseWriter, r *http.Request) {
	handleRawRequest(w, r)
}

// handleRawRequest returns the raw stats to http clients
func handleRawRequest(w http.ResponseWriter, r *http.Request) {
	// get all counters
	m := counter.getAll()

	// sort counters
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// send sorted counters to client
	for _, k := range keys {
		_, err := io.WriteString(w, fmt.Sprintln(k, m[k]))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// runHTTPServer runs the http server
func runHTTPServer(address string) {
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/raw", handleRawRequest)
	log.Fatal(http.ListenAndServe(address, nil))
}
