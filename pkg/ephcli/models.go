package ephcli

import "time"

const apiVersion string = "api/v1"

type ResponseError struct {
	Err     bool   `json:"error"`
	Message string `json:"msg"`
}

type File struct {
	Idfile          string    `json:"idfile"`
	FileName        string    `json:"filename"`
	Size            int64     `json:"size"`
	UpdateDateBegin time.Time `json:"update_date_egin"`
	UpdateDateEnd   time.Time `json:"update_date_end"`
	ExpirationDate  time.Time `json:"expiration_date"`
}

type Downloader interface {
	Download(uuidFileToDownload string, outputfile string) error
}

type Uploader interface {
	Upload(inputfile string) error
}

type FilesLister interface {
	List() error
}
