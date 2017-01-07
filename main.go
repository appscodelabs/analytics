package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/appscode/go/flags"
	flag "github.com/spf13/pflag"
)

// ref: http://stackoverflow.com/a/33301173
func main() {
	addr := flag.String("addr", ":60010", "host:port used by this server")
	caCertFile := flag.String("caCertFile", "", "File containing CA certificate")
	certFile := flag.String("certFile", "", "File container server TLS certificate")
	keyFile := flag.String("keyFile", "", "File containing server TLS private key")

	flags.InitFlags()
	flags.DumpAll()

	http.DefaultServeMux.HandleFunc("/index.json", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Println("Listening on", *addr)

	srv := &http.Server{
		Addr:         *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      http.DefaultServeMux,
	}
	if *caCertFile == "" && *certFile == "" && *keyFile == "" {
		log.Fatalln(srv.ListenAndServe())
	} else {
		/*
			Ref:
			 - https://blog.cloudflare.com/exposing-go-on-the-internet/
			 - http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/
			 - http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html
		*/
		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			SessionTicketsDisabled:   true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				// tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				// tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			ClientAuth: tls.VerifyClientCertIfGiven,
		}
		if *caCertFile != "" {
			caCert, err := ioutil.ReadFile(*caCertFile)
			if err != nil {
				log.Fatal(err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.ClientCAs = caCertPool
		}
		tlsConfig.BuildNameToCertificate()

		srv.TLSConfig = tlsConfig
		log.Fatalln(srv.ListenAndServeTLS(*certFile, *keyFile))
	}
}
