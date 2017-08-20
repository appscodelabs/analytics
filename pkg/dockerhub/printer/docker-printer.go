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
	results, err := dockerhub.GetRepoStats(org)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	// New Tab Writer
	w := getNewTabWriter()
	defer w.Flush()

	// New CSV writer
	var wCSV *csv.Writer
	if output != "" {
		wCSV, err = getNewCSVWriter(output, org)
		defer wCSV.Flush()
	}

	for _, repo := range results {
		//NAME,	PUL COUNT,	STAR COUNT,	LAST UPDATED
		str := fmt.Sprintf("%v/%v\t%v\t%v\t%v", repo.User, repo.Name, repo.PullCount, repo.StarCount, repo.LastUpdated.Format("2006-01-02 15:04"))
		fmt.Fprintln(w, str)
		if output != "" {
			wCSV.Write([]string{repo.Name, strconv.FormatInt(int64(repo.PullCount), 10), strconv.FormatInt(int64(repo.StarCount), 10), repo.LastUpdated.Format("2006-01-02 15:04")})
		}
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
