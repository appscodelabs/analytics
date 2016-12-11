package main

import (
	"fmt"
	"net"
	"net/http"
)

// ref: http://stackoverflow.com/a/33301173
func main() {
	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ip, port, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
				return
			}

			userIP := net.ParseIP(ip)
			if userIP == nil {
				fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
				return
			}

			// This will only be defined when site is accessed via non-anonymous proxy
			// and takes precedence over RemoteAddr
			// Header.Get is case-insensitive
			forward := r.Header.Get("X-Forwarded-For")

			fmt.Fprintf(w, "<p>IP: %s</p>", ip)
			fmt.Fprintf(w, "<p>Port: %s</p>", port)
			fmt.Fprintf(w, "<p>Forwarded for: %s</p>", forward)
		} else {

			http.Error(w, "Method not allowed", 405)
			return
		}
	})
	if err := http.ListenAndServe(":5050", nil); err == nil {
		fmt.Println(err)
	}
}
