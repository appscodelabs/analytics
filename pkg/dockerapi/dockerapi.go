package dockerapi

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

type DockerRepositories struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		User           string      `json:"user"`
		Name           string      `json:"name"`
		Namespace      string      `json:"namespace"`
		RepositoryType string      `json:"repository_type"`
		Status         int         `json:"status"`
		Description    string      `json:"description"`
		IsPrivate      bool        `json:"is_private"`
		IsAutomated    bool        `json:"is_automated"`
		CanEdit        bool        `json:"can_edit"`
		StarCount      int         `json:"star_count"`
		PullCount      int         `json:"pull_count"`
		LastUpdated    time.Time   `json:"last_updated"`
		BuildOnCloud   interface{} `json:"build_on_cloud"`
	} `json:"results"`
}

type DockerRepoLogs struct {
	User        string
	Name        string
	StarCount   int
	PullCount   int
	LastUpdated time.Time
}

const SpreadSheetIdKubeDB = "10OGrTJxEDox4VR15U7HPGRjiFpTNfhiUGMuL7BQJKJU"   //https://docs.google.com/spreadsheets/d/10OGrTJxEDox4VR15U7HPGRjiFpTNfhiUGMuL7BQJKJU
const SpreadSheetIdAppsCode = "18lNbYqiP4gsKBoLDoUw2ejOGWN3t3mUOyKGaKeye3kI" //https://docs.google.com/spreadsheets/d/18lNbYqiP4gsKBoLDoUw2ejOGWN3t3mUOyKGaKeye3kI

func getDockerLogs(urlLink string) (*DockerRepositories, error) {

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
	var record DockerRepositories

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return nil, err
	}
	return &record, nil
}

func updateSheet(dockLogs DockerRepoLogs, SpreadSheetId string) (string, error) {
	// Usase Limits: https://developers.google.com/sheets/api/limits
	// 100 Requeste per 100Seconds per User
	// 500 requests per 100 seconds per project
	// So, taking 5seconds before two requests
	time.Sleep(5 * time.Second)

	ctx := context.Background()

	//Reading Secret File from /.Credential
	b, err := getSecretFilePath()
	if err != nil {
		return "", errors.New("Unable to read client secret file").Err()
	}

	// If modifying these scopes, delete previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-api.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return "", errors.New("Unable to parse client secret file to config").Err()
	}
	client := spreadsheet.GetClient(ctx, config)
	srv, err := sheets.New(client)
	if err != nil {
		return "", errors.New("Unable to read client secret file: ").Err()
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

	if ok, err := srv.Spreadsheets.BatchUpdate(SpreadSheetId, batchRequest).Context(ctx).Do(); ok != nil {
		log.Println(dockLogs.Name, "sheet successfully created")
		values = append(values, []interface{}{"Last Updated", "Pull Count", "Star Count"})
	} else {
		log.Println(err)
		//Error because most probably sheet name already exists.
		//So, Do the rest of the work.
	}

	// Assign DockerRepoLogs to into values
	values = append(values, []interface{}{dockLogs.LastUpdated, dockLogs.PullCount, dockLogs.StarCount})

	rangeValue := dockLogs.Name + "!A:C"
	valueInputOption := "RAW"
	rb := &sheets.ValueRange{
		Values: values,
	}

	resp, err := srv.Spreadsheets.Values.Append(SpreadSheetId, rangeValue, rb).ValueInputOption(valueInputOption).Context(ctx).Do()

	if err != nil {
		return "", errors.FromErr(err).Err()
	}
	return fmt.Sprintf("Successful [%v] row insertion in [%v] for [%v/%v]", resp.Updates.UpdatedRows, resp.SpreadsheetId, dockLogs.User, dockLogs.Name), nil
}

func processAndUpdate(spreadSheetId string, link string) {
	dockerResp, err := getDockerLogs(link)
	if err != nil {
		log.Fatalf("Unable to retrieve Docker log.s %v", err)
	}

	for _, c := range dockerResp.Results {
		respStr, err := updateSheet(DockerRepoLogs{
			Name:        c.Name,
			User:        c.User,
			StarCount:   c.StarCount,
			PullCount:   c.PullCount,
			LastUpdated: time.Now(),
		}, spreadSheetId)

		if err != nil {
			log.Println("Error while updating ", c.Name, err)
		} else {
			log.Println(respStr)
		}
	}

	log.Println(dockerResp.Next)
	for dockerResp.Next != nil {
		dockerResp, err = getDockerLogs(*dockerResp.Next)

		if err != nil {
			log.Fatalf("Unable to retrieve Docker log.s %v", err)
		}
		for _, c := range dockerResp.Results {
			respStr, err := updateSheet(DockerRepoLogs{
				Name:        c.Name,
				User:        c.User,
				StarCount:   c.StarCount,
				PullCount:   c.PullCount,
				LastUpdated: time.Now(),
			}, spreadSheetId)
			if err != nil {
				log.Println("Error while updating ", c.Name, err)
			} else {
				log.Println(respStr)
			}
		}
		log.Println(dockerResp.Next)
	}
}

// File Path: ~/.credentials/client_secret_spreadsheet.json
func getSecretFilePath() ([]byte, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	fileDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(fileDir, 0700)
	return ioutil.ReadFile(filepath.Join(fileDir, url.QueryEscape("client_secret_spreadsheet.json")))
}

func deleteSheet(SpreadSheetId string) {
	ctx := context.Background()

	//Reading Secret File from /.Credential
	b, err := getSecretFilePath()
	if err != nil {
		log.Println("Unable to read client secret file", err)
	}

	// If modifying these scopes, delete previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-api.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Println("Unable to parse client secret file to config")
	}
	client := spreadsheet.GetClient(ctx, config)
	srv, err := sheets.New(client)
	if err != nil {
		log.Println("Unable to read client secret file: ")
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

		if ok, err := srv.Spreadsheets.BatchUpdate(SpreadSheetId, batchRequest).Context(ctx).Do(); ok != nil {
			log.Println("sheet successfully Deleted")
		} else {
			log.Println(err)
			//Error because most probably sheet name already exists.
			//So, Do the rest of the work.
		}
	}
}

func DockerAnalytics() {

	processAndUpdate(SpreadSheetIdKubeDB, "https://hub.docker.com/v2/repositories/kubedb/?page_size=100")
	processAndUpdate(SpreadSheetIdAppsCode, "https://hub.docker.com/v2/repositories/appscode/?page_size=100")

	// To delete all the sheets of a Spreadsheet
	//deleteSheet(SpreadSheetIdKubeDB)
	//deleteSheet(SpreadSheetIdAppsCode)
}
