package main

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
	"strings"
	"sync"

	filepath_x "github.com/a8uhnf/go-utility/filepath"
	"github.com/fsnotify/fsnotify"
	"github.com/ghodss/yaml"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// ConfigFileName is name of config file.
const ConfigFileName = "config.yaml"

// Config contains config of spreadsheet
type Config struct {
	SpreadsheetID string `json:"spreadsheetID"`
	SheetName     string `json:"sheetName"`
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	fmt.Println("Hello World-------", usr)
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
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
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getConfig() (*Config, error) {
	p, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(p, ConfigFileName)
	y, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(y, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	ctx := context.Background()
	/*b, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), "credentials/google-spreadsheet/client_secret.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := cfg.SpreadsheetID
	readRange := cfg.SheetName + "!A1:A" */
	resp, err := GetSpreadsheetData(ctx) // srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
	fmt.Println("------------------", len(resp.Values))
	if len(resp.Values) > 0 {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			// fmt.Printf("%s, %s\n", row[0])
			if len(row) > 0 {
				fmt.Println(row)
			}
		}
	} else {
		fmt.Print("No data found.")
	}
	err = DownloadWatcher(ctx)
	if err != nil {
		panic(err)
	}
}

// DownloadWatcher watches download folder.
func DownloadWatcher(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		appender := &SpreadsheetAppender{}
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("-------------------")
				log.Println("event operations:- ", event.Op)
				log.Println("filename:- ", event.Name)
				if strings.HasSuffix(event.Name, ".torrent") && event.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Println("Hello Find the files.....")
					svc, err := getSpreadsheetService(ctx)
					if err != nil {
						panic(err)
					}
					err = appender.AppendSpreadSheet(ctx, svc, filepath_x.Filename(event.Name))
					if err != nil {
						panic(err)
					}
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	path := filepath.Join(os.Getenv("HOME"), "Downloads")
	if _, ok := os.Stat(path); ok != nil {
		return ok
	}
	err = watcher.Add(path)
	if err != nil {
		return err
	}
	<-done
	return nil
}

func GetSpreadsheetData(ctx context.Context) (*sheets.ValueRange, error) {
	srv, err := getSpreadsheetService(ctx)
	if err != nil {
		return nil, err
	}
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := cfg.SpreadsheetID
	readRange := cfg.SheetName + "!A1:A"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getSpreadsheetService(ctx context.Context) (*sheets.Service, error) {
	b, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), "credentials/google-spreadsheet/client_secret.json"))
	if err != nil {
		return nil, err
	}
	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

// SpreadsheetAppender holds values for append to spreadsheet
type SpreadsheetAppender struct {
	lock sync.Mutex
}

// AppendSpreadSheet appends to spreadsheet.
func (s *SpreadsheetAppender) AppendSpreadSheet(ctx context.Context, sheetsService *sheets.Service, name string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	// The ID of the spreadsheet to update.
	spreadsheetId := cfg.SpreadsheetID // TODO: Update placeholder value.
	l, err := GetSpreadsheetData(ctx)
	if err != nil {
		return err
	}
	// The A1 notation of a range to search for a logical table of data.
	// Values will be appended after the last row of the table.
	range2 := fmt.Sprintf("A%d:A%d", len(l.Values), len(l.Values)) // TODO: Update placeholder value.
	fmt.Println("Range:= ", range2)
	// How the input data should be interpreted.
	valueInputOption := "USER_ENTERED" // TODO: Update placeholder value.

	// How the input data should be inserted.
	insertDataOption := "OVERWRITE" // TODO: Update placeholder value.

	test := [][]interface{}{}
	test = append(test, []interface{}{name})
	fmt.Println(test)
	rb := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          range2,
		Values:         test,
	}

	resp, err := sheetsService.Spreadsheets.Values.Append(spreadsheetId, range2, rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
		return err
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("%#v\n", resp)

	return nil
}
