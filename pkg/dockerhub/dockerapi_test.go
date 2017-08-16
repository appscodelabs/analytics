package dockerhub

import (
	"testing"

	"github.com/appscode/log"
)

func TestProcessAndUpdate(t *testing.T) {
	err := refresh("10OGrTJxEDox4VR15U7HPGRjiFpTNfhiUGMuL7BQJKJU", "https://hub.docker.com/v2/repositories/kubedb/?page_size=100")
	if err != nil {
		log.Fatalln(err)
	}
}
