package ephcli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type downloader struct {
	HttpClient  *http.Client
	token       string
	endpoint    string
	bar         *mpb.Bar
	progressBar *mpb.Progress
}

// NewDownloader creates a new Downloader
func NewDownloader(endpoint string, token string) Downloader {
	return &downloader{
		HttpClient: &http.Client{},
		token:      token,
		endpoint:   fmt.Sprintf("%s/api/v1/download", endpoint),
	}
}

// Download downloads a file from the server
// and saves it to the outputfile
// If the outputfile is empty, the file will be saved to the current directory
// with the same name as the file on the server (retrieving the name from the Content-Disposition header)
func (d *downloader) Download(uuidFileToDownload string, outputfile string) error {
	var filename string
	url := fmt.Sprintf("%s/%s", d.endpoint, uuidFileToDownload)
	// prepare request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))
	resp, err := d.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return parseError(resp)
	}
	filename = d.getFileName(resp, outputfile)
	totalSize := resp.ContentLength
	d.initProgressBar(filename)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			stat, err := f.Stat()
			if err != nil {
				fmt.Println(err)
			}
			d.bar.SetCurrent(int64(float64(stat.Size()) / float64(totalSize) * 100.0))
			time.Sleep(1 * time.Second)
			if stat.Size() == totalSize {
				break
			}
		}
		wg.Done()
	}()
	io.Copy(f, resp.Body)
	wg.Wait()
	d.bar.Wait()
	defer resp.Body.Close()
	return nil
}

// getFileName returns outputFileName if not empty
// If empty, try to retrieve the filename from the Content-Disposition header
// If not present, return the last part of the URL
func (d *downloader) getFileName(resp *http.Response, outputfileName string) string {
	var filename string
	if outputfileName != "" {
		return outputfileName
	}
	contentDisposition := resp.Header.Get("Content-Disposition")
	// extract filename from contentDisposition
	contentDispositionSplitted := strings.Split(contentDisposition, "=")
	if len(contentDispositionSplitted) < 2 {
		filename = filepath.Base(resp.Request.URL.Path)
	} else {
		filename = contentDispositionSplitted[1]
		// remove double quotes at the begin and the end from filename with regexp
		filename = strings.TrimLeft(strings.TrimRight(filename, "\""), "\"")
	}
	return filename
}

func (d *downloader) initProgressBar(fileToDownload string) {
	// initialize progress container, with custom width
	d.progressBar = mpb.New(mpb.WithWidth(64))
	total := 100
	name := filepath.Base(fileToDownload)
	// create a single bar, which will inherit container's width
	d.bar = d.progressBar.New(int64(total),
		// BarFillerBuilder with custom style
		mpb.BarStyle(),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(name, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			// replace ETA decorator with "done" message, OnComplete event
			decor.OnComplete(decor.Elapsed(decor.ET_STYLE_GO), "done"),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)
}
