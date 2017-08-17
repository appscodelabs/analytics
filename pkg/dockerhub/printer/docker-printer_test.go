package printer_test

import (
	"log"
	"testing"

	"github.com/appscode/analytics/pkg/dockerhub/printer"
)

func TestDockerPrinter(t *testing.T) {
	err := printer.DockerPrinter([]string{"appscode", "kubedb"}, "/Desktop/Temporary")
	if err != nil {
		log.Println(err)
	}
}
