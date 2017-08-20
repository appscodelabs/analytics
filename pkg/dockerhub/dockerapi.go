package dockerhub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/appscode/analytics/pkg/spreadsheet"
	"github.com/appscode/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

type RepoStats struct {
	User        string    `json:"user"`
	Name        string    `json:"name"`
	Namespace   string    `json:"namespace"`
	StarCount   int       `json:"star_count"`
	PullCount   int       `json:"pull_count"`
	LastUpdated time.Time `json:"last_updated"`
}

type OrgStats struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []RepoStats `json:"results"`
}

func getDockerLogs(urlLink string) (*OrgStats, error) {
	// Build the request
	req, err := http.NewRequest("GET", urlLink, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record OrgStats

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return nil, err
	}
	return &record, nil
}

func updateSheet(srv *sheets.Service, spreadSheetID string, repo RepoStats, blank bool) error {
	// Usase Limits: https://developers.google.com/sheets/api/limits
	// 100 Requeste per 100Seconds per User
	// 500 requests per 100 seconds per project
	// So, taking 5seconds before two requests
	time.Sleep(5 * time.Second)
	ctx := context.Background()
	var values [][]interface{}
	// Assign RepoStats to into values
	if blank {
		values = append(values, []interface{}{"Timestamp", "Pull Count", "Star Count"})
	}
	values = append(values, []interface{}{time.Now().Format("2006-01-02 15:04:05 Z07:00"), repo.PullCount, repo.StarCount})
	rangeValue := repo.Name + "!A:C"
	valueInputOption := "RAW"
	rb := &sheets.ValueRange{
		Values: values,
	}
	resp, err := srv.Spreadsheets.Values.Append(spreadSheetID, rangeValue, rb).ValueInputOption(valueInputOption).Context(ctx).Do()
	if err != nil {
		return errors.FromErr(err).Err()
	}
	log.Printf("Successful [%v] row insertion in [%v] for [%v/%v]\n", resp.Updates.UpdatedRows, resp.SpreadsheetId, repo.User, repo.Name)
	return nil
}

func createSheet(srv *sheets.Service, spreadSheetID string, repo RepoStats) error {
	ctx := context.Background()
	requests := []*sheets.Request{}
	requests = append(requests, &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: repo.Name,
				GridProperties: &sheets.GridProperties{
					ColumnCount: 3,
				},
			},
		},
	})
	batchRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}
	if _, err := srv.Spreadsheets.BatchUpdate(spreadSheetID, batchRequest).Context(ctx).Do(); err != nil {
		return err
	}
	return nil
}

func GetRepoStats(org string) ([]RepoStats, error) {
	link := fmt.Sprintf("https://hub.docker.com/v2/repositories/%v/?page_size=50", org)
	dockerResp, err := getDockerLogs(link)
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}
	result := make([]RepoStats, 0)

	first := true
	for first || dockerResp.Next != nil {
		if !first {
			dockerResp, err = getDockerLogs(*dockerResp.Next)
			if err != nil {
				return nil, errors.FromErr(err).Err()
			}
		}
		result = append(result, dockerResp.Results...)
		first = false
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

func sheetAndStats(sheetid, org string) error {
	result, err := GetRepoStats(org)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	srv, err := spreadsheet.GetNewSheetService()
	if err != nil {
		return errors.FromErr(err).Err()
	}
	for _, c := range result {
		blank := false
		if err := createSheet(srv, sheetid, c); err != nil {
			log.Println(err)
			//Error because most probably sheet name already exists.
			//So, Do the rest of the work.
		} else {
			log.Println(c.Name, "sheet successfully created")
			blank = true
		}

		if err = updateSheet(srv, sheetid, c, blank); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func CollectAnalytics(dockerOrgs map[string]string) error {
	for org, sheetID := range dockerOrgs {
		err := sheetAndStats(sheetID, org)
		if err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}
