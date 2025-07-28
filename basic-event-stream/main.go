package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	flag.Parse()
	addr := flag.String("addr", ":8080", "HTTP network address")
	http.HandleFunc("/", home)
	http.HandleFunc("/events", events)
	http.ListenAndServe(*addr, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")

	token := []string{"this", "is", "a", "test", "event"}

	for _, t := range token {
		content := fmt.Sprintf("data: %s\n\n", t)
		w.Write([]byte(content))
		w.(http.Flusher).Flush()

		time.Sleep(time.Millisecond * 420)
	}
}
