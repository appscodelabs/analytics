package pixels

import (
	"net/http"
	"strings"

	"github.com/appscode/errors"
	"github.com/appscode/pat"
	"github.com/jpillora/go-ogle-analytics"
)

// Tracking pixel. Ref: https://product.reverb.com/build-a-protocol-buffer-powered-tracking-pixel-in-go-76f2ca5c26e2
var GIF = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

func ImageHits(w http.ResponseWriter, req *http.Request) {
	params, found := pat.FromContext(req.Context())
	if !found {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}
	path := strings.SplitN(req.URL.Path, "/", 4)[3]
	err := sendPageView(params.Get(":trackingcode"), params.Get(":host"), req.RemoteAddr, req.UserAgent(), path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/gif")
	w.Write(GIF)
}

func sendPageView(trackingcode, host, ip, userAgent, path string) error {
	client, err := ga.NewClient(trackingcode)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	client.DocumentHostName(host)
	client.IPOverride(ip)
	client.UserAgentOverride(userAgent)
	client.DocumentPath(path)

	err = client.Send(ga.NewPageview())
	return err
}
