package pixels

import (
	"net/http"
	"strings"

	"github.com/appscode/analytics/pkg/analytics"
	"github.com/appscode/pat"
)

// tracking code = "UA-62096468-19"

// Tracking pixel
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
	analytics.SendPageView(params.Get(":trackingcode"), params.Get(":host"), req.RemoteAddr, req.UserAgent(), path)
	w.Write(GIF)
}
