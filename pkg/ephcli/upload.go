package ephcli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/imroc/req/v3"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// https://req.cool/docs/tutorial/download/

type uploader struct {
	HttpClient  *http.Client
	token       string
	endpoint    string
	bar         *mpb.Bar
	progressBar *mpb.Progress
}

func NewUploader(endpoint string, token string) Uploader {
	return &uploader{
		HttpClient: &http.Client{},
		token:      token,
		endpoint:   fmt.Sprintf("%s/api/v1/upload", endpoint),
	}
}

func (d *uploader) initProgressBar(fileToUpload string) {
	// initialize progress container, with custom width
	d.progressBar = mpb.New(mpb.WithWidth(64))
	total := 100
	name := filepath.Base(fileToUpload)
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

func (d *uploader) callback(info req.UploadInfo) {
	d.bar.SetCurrent(int64(float64(info.UploadedSize) / float64(info.FileSize) * 100.0))
}

func (d *uploader) Upload(fileToUpload string) error {
	d.initProgressBar(fileToUpload)
	client := req.C().SetCommonBearerAuthToken(d.token).SetTimeout(0)
	resp, err := client.R().
		SetFile("uploadfile", fileToUpload).
		SetUploadCallback(d.callback).Post(d.endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return d.parseError(resp)
	}
	d.progressBar.Wait()
	return nil
}

// parseError parses the error from the server
func (d *uploader) parseError(resp *req.Response) error {
	var jsonResponse ResponseError
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &jsonResponse)
	if err != nil {
		return err
	}
	return fmt.Errorf("status not ok %d: %s", resp.StatusCode, jsonResponse.Message)
}
