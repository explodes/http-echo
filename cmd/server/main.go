package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	flagBind = flag.String("bind", ":8087", "bind address")
)

func init() {
	flag.Parse()
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	query := map[string][]string(url.Query())
	headers := map[string][]string(r.Header)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	writeString(w, "Path: %s", url)

	if len(query) > 0 {
		writeString(w, "\nQuery Paramters")
		writeMapStringSliceString(w, query)
	}

	if len(headers) > 0 {
		writeString(w, "\nHeaders")
		writeMapStringSliceString(w, headers)
	}
}

func writeMapStringSliceString(w http.ResponseWriter, data map[string][]string) {
	for k, v := range data {
		if len(v) == 0 {
			writeString(w, " - %s: <empty>", k)
		} else if len(v) == 1 {
			writeString(w, " - %s: %s", k, v[0])
		} else {
			for index, value := range v {
				writeString(w, " - %s #%d: %s", k, index+1, value)
			}
		}
	}
}

func writeString(w http.ResponseWriter, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	w.Write([]byte(msg))
	w.Write([]byte("\n"))
}

func onError(w http.ResponseWriter, cause string, err error) {
	log.Printf("%s: %v", cause, err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func main() {
	log.Printf("Serving on %s", *flagBind)
	http.HandleFunc("/", echoHandler)
	log.Println(http.ListenAndServe(*flagBind, nil))
}
