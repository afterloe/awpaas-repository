package warehouse

import (
	"../../exceptions"
	"fmt"
	"reflect"
)

type warehouse struct {
	rev string
	Id string `json:"id"`
	Status bool `json:"status"`
	ModifyTime int64 `json:"modifyTime"`
	UploadTime int64 `json:"uploadTime"`
	Name string `json:"name"`
	Group string `json:"group"`
	Remarks string `json:"remarks"`
	Version string `json:"version"`
	Fid string `json:"fid"`
	PackInfo map[string]interface{} `json:"packInfo"`
	Cmd cmd `json:"cmd"`
}

var registryType = [3]string{}

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