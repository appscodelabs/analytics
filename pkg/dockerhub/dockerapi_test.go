package dockerhub

import (
	"testing"

	"github.com/appscode/log"
)

func TestRefreshStats(t *testing.T) {
	err := refreshStats("10OGrTJxEDox4VR15U7HPGRjiFpTNfhiUGMuL7BQJKJU", "https://hub.docker.com/v2/repositories/kubedb/?page_size=100")
	if err != nil {
		log.Fatalln(err)
	}
}
