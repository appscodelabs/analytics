package spreadsheet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/appscode/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var Client_secret_path string

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func GetClient(ctx context.Context, config *oauth2.Config) (*http.Client,error) {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		log.Println(err)
		// Browse the token and saves in file directory
		tok,err = getTokenFromWeb(config)
		if err != nil{
			return nil,errors.FromErr(err).Err()
		}
		err = saveToken(cacheFile, tok)
		if err != nil{
			return nil,errors.FromErr(err).Err()
		}
	}
	return config.Client(ctx, tok),nil
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token,error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil,errors.FromErr(err).Err()
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil,errors.FromErr(err).Err()
	}
	return tok,nil
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	tokenCacheDir := filepath.Join(Client_secret_path)
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheet_access_token.json")), nil
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

func GetNewSheetService() (*sheets.Service, error) {
	ctx := context.Background()
	b, err := getClientSecret()
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}
	// If modifying these scopes, delete previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-api.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, errors.New("Unable to parse client secret file to config").Err()
	}
	client,err := GetClient(ctx, config)
	if err !=nil{
		return nil,errors.FromErr(err).Err()
	}
	srv, err := sheets.New(client)
	return srv, err
}

// File Path: /Client_Secret_Path
func getClientSecret() ([]byte, error) {
	fileDir := filepath.Join(Client_secret_path)
	os.MkdirAll(fileDir, 0700)
	return ioutil.ReadFile(filepath.Join(fileDir, url.QueryEscape("client_secret_spreadsheet.json")))
}
