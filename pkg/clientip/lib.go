package clientip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func WhoAmI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, fmt.Sprintf("userip: %q is not IP:port", r.RemoteAddr), 500)
			return
		}

		userIP := net.ParseIP(ip)
		if userIP == nil {
			http.Error(w, fmt.Sprintf("userip: %q is not IP:port", r.RemoteAddr), 500)
			return
		}

		// This will only be defined when site is accessed via non-anonymous proxy
		// and takes precedence over RemoteAddr
		forward := r.Header.Get("X-Forwarded-For")

		data := map[string]interface{}{
			"ip":            ip,
			"forwarded_for": forward,
		}
		if strings.EqualFold(r.URL.Query().Get("include_headers"), "true") {
			data["request_headers"] = r.Header
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	} else {
		http.Error(w, "Method not allowed", 405)
		return
	}
}
