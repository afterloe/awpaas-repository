package borderSystem

import (
	"fmt"
	"time"
)

var (
	timeFormat string
)

func init() {
	timeFormat = "2006-01-02 - 15:04:05"
}

type fsFile struct {
	Id string `json:"id"`
	Name string `json:"name"`
	SavePath string `json:"savePath"`
	ContentType string `json:"contentType"`
	Key string `json:"key"`
	UploadTime int64 `json:"uploadTime"`
	Size int64 `json:"size"`
	Status bool `json:"status"`
	ModifyTime int64 `json:"modifyTime"`
	rev string
}

func (this *fsFile) String() string {
	return fmt.Sprintf("name: %s savePaht: %s contentType: %s key: %s, uploadTime: %s, size: %d, status %v",
		this.Name, this.SavePath, this.ContentType, this.Key, time.Unix(this.UploadTime, 0).Format(timeFormat),
		this.Size, this.Status)
}