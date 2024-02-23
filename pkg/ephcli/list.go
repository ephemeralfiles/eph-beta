package ephcli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type lister struct {
	HttpClient *http.Client
	token      string
	endpoint   string
}

func NewLister(endpoint string, token string) FilesLister {
	return &lister{
		HttpClient: &http.Client{},
		token:      token,
		endpoint:   fmt.Sprintf("%s/%s/files", endpoint, apiVersion),
	}
}

func (l *lister) List() error {
	req, err := http.NewRequest("GET", l.endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", l.token))
	resp, err := l.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseError(resp)
	}

	var files []File
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Fprintf(os.Stdout, "%37s %20s %10d %20s\n", file.Idfile, file.FileName, file.Size, file.ExpirationDate.Format("2006-01-02 15:04:05"))
	}
	return nil
}
