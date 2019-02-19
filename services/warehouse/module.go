package warehouse

import (
	"../../exceptions"
	"fmt"
	"reflect"
	"time"
)

type warehouse struct {
	rev string
	Id int64 `json:"id"`
	Status bool `json:"status"`
	ModifyTime int64 `json:"modifyTime"`
	UploadTime int64 `json:"uploadTime"`
	Name string `json:"name"`
	Group string `json:"group"`
	Remarks string `json:"remarks"`
	Version string `json:"version"`
	PackInfo fsFile `json:"packInfo"`
	Cmd cmd `json:"cmd"`
}

var registryType = [4]string{}

type cmd struct {
	RegistryType string `json:"registryType"`
	Content []string `json:"content"`
	LastReport string `json:"lastReport"`
	LastCiTime int64 `json:"lastCiTime"`
}

func (this *cmd) String() string {
	return fmt.Sprintf("{'fileType': '%s', 'content': '%v'}", this.RegistryType, this.Content)
}

func (this *warehouse) String() string {
	return fmt.Sprintf("{'name': '%s', 'id': '%s', 'rev': '%s'}", this.Name, this.Id, this.rev)
}

/**
	参数检测
*/
func (this *warehouse) Check(args ...string) error {
	value := reflect.ValueOf(*this)
	for _, arg := range args {
		v := value.FieldByName(arg)
		if !v.IsValid() {
			break
		}
		if "" == v.Interface() {
			return &exceptions.Error{Msg: "lack param " + arg, Code: 400}
		}
	}

	return nil
}

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
}

func (this *fsFile) String() string {
	return fmt.Sprintf("name: %s savePaht: %s contentType: %s key: %s, uploadTime: %s, size: %d, status %v",
		this.Name, this.SavePath, this.ContentType, this.Key, time.Unix(this.UploadTime, 0).Format(timeFormat),
		this.Size, this.Status)
}