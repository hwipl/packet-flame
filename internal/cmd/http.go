package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sort"
	"strings"
)

// handleRequest returns the flamegraph to http clients
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// if no flamegraph tool is configured, just return raw data
	if flamegraphExe == "" {
		handleRawRequest(w, r)
		return
	}

	// generate input string for flamegraph tool
	input := ""
	for l, c := range counter.getAll() {
		input += fmt.Sprintln(l, c)
	}

	// generate flamegraph
	cmd := exec.Command(flamegraphExe)
	if flamegraphArgs != "" {
		args := strings.Split(flamegraphArgs, " ")
		cmd = exec.Command(flamegraphExe, args...)
	}
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = w
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
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
