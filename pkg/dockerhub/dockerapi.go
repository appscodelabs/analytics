package dockerhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/appscode/analytics/pkg/spreadsheet"
	"github.com/appscode/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type OrgStats struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		User        string    `json:"user"`
		Name        string    `json:"name"`
		Namespace   string    `json:"namespace"`
		StarCount   int       `json:"star_count"`
		PullCount   int       `json:"pull_count"`
		LastUpdated time.Time `json:"last_updated"`
	} `json:"results"`
}

type DockerRepoLogs struct {
	User      string
	Name      string
	StarCount int
	PullCount int
	Timestamp time.Time
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

func updateSheet(dockLogs DockerRepoLogs, SpreadSheetId string) error {
	// Usase Limits: https://developers.google.com/sheets/api/limits
	// 100 Requeste per 100Seconds per User
	// 500 requests per 100 seconds per project
	// So, taking 5seconds before two requests
	time.Sleep(5 * time.Second)

	ctx := context.Background()
	b, err := getClientSecret()
	if err != nil {
		return errors.New("Unable to read client secret file").Err()
	}

	// If modifying these scopes, delete previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-api.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return errors.New("Unable to parse client secret file to config").Err()
	}
	client := spreadsheet.GetClient(ctx, config)
	srv, err := sheets.New(client)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	var values [][]interface{}
	// Create sheet if not exists.
	requests := []*sheets.Request{}
	requests = append(requests, &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: dockLogs.Name,
				GridProperties: &sheets.GridProperties{
					ColumnCount: 3,
					//because only 2000000 cells are allowed!
					//ref:
				},
			},
		},
	})
	batchRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}
	if _, err := srv.Spreadsheets.BatchUpdate(SpreadSheetId, batchRequest).Context(ctx).Do(); err == nil {
		log.Println(dockLogs.Name, "sheet successfully created")
		values = append(values, []interface{}{"Timestamp", "Pull Count", "Star Count"})
	} else {
		log.Println(err)
		//Error because most probably sheet name already exists.
		//So, Do the rest of the work.
	}

	// Assign DockerRepoLogs to into values
	values = append(values, []interface{}{dockLogs.Timestamp, dockLogs.PullCount, dockLogs.StarCount})
	rangeValue := dockLogs.Name + "!A:C"
	valueInputOption := "RAW"
	rb := &sheets.ValueRange{
		Values: values,
	}
	resp, err := srv.Spreadsheets.Values.Append(SpreadSheetId, rangeValue, rb).ValueInputOption(valueInputOption).Context(ctx).Do()
	if err != nil {
		return errors.FromErr(err).Err()
	}
	log.Printf("Successful [%v] row insertion in [%v] for [%v/%v]\n", resp.Updates.UpdatedRows, resp.SpreadsheetId, dockLogs.User, dockLogs.Name)
	return nil
}

func refresh(spreadSheetId string, link string) error {
	dockerResp, err := getDockerLogs(link)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	for _, c := range dockerResp.Results {
		err := updateSheet(DockerRepoLogs{
			Name:      c.Name,
			User:      c.User,
			StarCount: c.StarCount,
			PullCount: c.PullCount,
			Timestamp: time.Now(),
		}, spreadSheetId)
		if err != nil {
			return errors.FromErr(err).Err()
		}
	}
	for dockerResp.Next != nil {
		dockerResp, err = getDockerLogs(*dockerResp.Next)
		if err != nil {
			return errors.FromErr(err).Err()
		}
		for _, c := range dockerResp.Results {
			err := updateSheet(DockerRepoLogs{
				Name:      c.Name,
				User:      c.User,
				StarCount: c.StarCount,
				PullCount: c.PullCount,
				Timestamp: time.Now(),
			}, spreadSheetId)
			if err != nil {
				return errors.FromErr(err).Err()
			}
		}
	}
	return nil
}

// File Path: ~/.credentials/client_secret_spreadsheet.json
func getClientSecret() ([]byte, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	fileDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(fileDir, 0700)
	return ioutil.ReadFile(filepath.Join(fileDir, url.QueryEscape("client_secret_spreadsheet.json")))
}

func deleteSheet(SpreadSheetId string) error {
	ctx := context.Background()
	b, err := getClientSecret()
	if err != nil {
		return errors.FromErr(err).Err()
	}

	// If modifying these scopes, delete previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-api.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return errors.FromErr(err).Err()
	}
	client := spreadsheet.GetClient(ctx, config)
	srv, err := sheets.New(client)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	docResp, err := srv.Spreadsheets.Get(SpreadSheetId).Context(ctx).Do()
	for _, c := range docResp.Sheets {
		// Delete sheets if exists
		requests := []*sheets.Request{}
		requests = append(requests, &sheets.Request{
			DeleteSheet: &sheets.DeleteSheetRequest{
				SheetId: c.Properties.SheetId,
			},
		})
		batchRequest := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}
		if _, err := srv.Spreadsheets.BatchUpdate(SpreadSheetId, batchRequest).Context(ctx).Do(); err == nil {
			log.Println("sheet successfully Deleted")
		} else {
			return errors.FromErr(err).Err()
			//Error because most probably sheet name already exists.
			//So, Do the rest of the work.
		}
	}
	return nil
}

func CollectAnalytics(dockerOrgs map[string]string) error {
	for org, sheetID := range dockerOrgs {
		err := refresh(sheetID, fmt.Sprintf("https://hub.docker.com/v2/repositories/%v/?page_size=50", org))
		if err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}
