package printer

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appscode/analytics/pkg/dockerhub"
	"github.com/appscode/errors"
)

func DockerPrinter(strArr []string) error {
	for _, org := range strArr {
		err := processAndPrint(fmt.Sprintf("https://hub.docker.com/v2/repositories/%v/?page_size=50", org))
		if err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func processAndPrint(link string) error {
	dockerResp, err := dockerhub.GetDockerLogs(link)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	w := getNewTabWriter()
	first := true
	for first || dockerResp.Next != nil {
		if !first {
			dockerResp, err = dockerhub.GetDockerLogs(*dockerResp.Next)
			if err != nil {
				return errors.FromErr(err).Err()
			}
		}
		for _, c := range dockerResp.Results {
			//NAME,	PUL COUNT,	STAR COUNT,	LAST UPDATED
			str := fmt.Sprintf("%v/%v\t%v\t%v\t%v", c.User, c.Name, c.PullCount, c.StarCount, c.LastUpdated)
			fmt.Fprintln(w, str)
		}
		first = false
	}
	fmt.Fprintln(w, "\n")
	w.Flush()

	return nil
}

func getNewTabWriter() *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tPULL COUNT\tSTAR COUNT\tLAST UPDATED")
	return w
}
