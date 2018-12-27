package test

import (
	"testing"
	"reflect"
	"encoding/json"
	"fmt"
)

type Module struct {
	Name, Group, Remarks, Version, Fid string
	_id string
}

func (this *Module) check(args ...string) {
	value := reflect.ValueOf(*this)
	for _, arg := range args {
		v := value.FieldByName(arg)
		if !v.IsValid() {
			break
		}
		if "" == v.Interface() {
			fmt.Println("11111111111 null")
			return
		}
	}
}

func Test_CheckModule(t *testing.T) {
	module := &Module {
		Name: "name",
		Group: "group",
		Remarks: "remarks",
		Version : "version",
		Fid: "fid",
		_id: "12",
	}
	buf, err := json.Marshal(module)
	if nil != err {
		t.Error(err)
	}
	t.Log(string(buf))
	module.check("Name", "Fid")
}
