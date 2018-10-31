package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().UTC().Format("2006-01-02 15:04:05.999")

	fmt.Printf("%v %v\t%v\n", t, r.Method, r.RequestURI)

	resp := make([]string, len(r.Header)+1)
	i := 1
	resp[0] = t
	for k, v := range r.Header {
		str := k + ":" + strings.Join(v, "\n\t")
		resp[i] = str
		i++
	}
	str := strings.Join(resp, "\n")

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(str))

	if err != nil {
		fmt.Printf("Can't send response: %v\n\t%v\n", string(str), err)
		http.Error(w, "Error: "+err.Error(), 501)
	}
}

func main() {
	host := flag.String("host", "localhost", "Host to listen")
	port := flag.Int("port", 8080, "Port to listen")

	addr := fmt.Sprintf("%v:%d", *host, *port)

	fmt.Printf("Listening :%v\n", addr)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		panic(err)
	}
}
