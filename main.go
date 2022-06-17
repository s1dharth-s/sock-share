package main

import (
	"net/http"

	"github.com/s1dharth-s/sock-share/sockshare"
)

func main() {
	h := sockshare.NewHub()
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		sockshare.HandleSocket(h, w, r)
	})
	go h.Run()
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe("localhost:8000", nil)
}
