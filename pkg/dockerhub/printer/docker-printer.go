package printer

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/appscode/analytics/pkg/dockerhub"
	"github.com/appscode/errors"
)

func DockerPrinter(strArr []string, output string) error {

	for _, org := range strArr {
		err := processAndPrint(output, org)
		if err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func processAndPrint(output string, org string) error {
	link := fmt.Sprintf("https://hub.docker.com/v2/repositories/%v/?page_size=50", org)
	dockerResp, err := dockerhub.GetDockerLogs(link)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	// New Tab Writer
	w := getNewTabWriter()
	defer w.Flush()

	// CSV writer
	var wCSV *csv.Writer
	if output != "" {
		wCSV, err = getNewCSVWriter(output, org)
		defer wCSV.Flush()
	}

	// Fetch dockerhub logs and write to writer
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
			if output != "" {
				wCSV.Write([]string{c.Name, strconv.FormatInt(int64(c.PullCount), 10), strconv.FormatInt(int64(c.StarCount), 10), c.LastUpdated.String()})
			}
		}
		first = false
	}
	fmt.Fprintln(w, "\n")
	return nil
}

func getNewTabWriter() *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tPULL COUNT\tSTAR COUNT\tLAST UPDATED")
	return w
}

func getNewCSVWriter(output string, org string) (*csv.Writer, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	// create Directory in output filepath
	fileDir := filepath.Join(usr.HomeDir, output)
	if err = os.MkdirAll(fileDir, 0700); err != nil {
		return nil, err
	}

	// create org-timestamp.csv file
	TimestampFormat := "20060102T1504"
	name := fmt.Sprintf("%v-%v.csv", org, time.Now().UTC().Format(TimestampFormat))
	fileDir = filepath.Join(fileDir, name)
	file, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}

	// CSV writer
	wCSV := csv.NewWriter(file)
	wCSV.Write([]string{"NAME", "PULL COUNT", "STAR COUNT", "LAST UPDATED"})
	return wCSV, nil
}
