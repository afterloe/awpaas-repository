package test

import (
	"testing"
	"strings"
	"fmt"
)

func Test_execArg(t *testing.T) {
	filename := "risk-point-0.0.5.tar.gz"
	filename = strings.ToLower(filename)
	names := strings.Split(filename,".tar.gz")
	t.Log(names)
	if "" != names[1] {
		t.Error(filename)
	}
	it := strings.Split(names[0], "-")
	name := strings.Join(it[:len(it) - 1], "-")
	version := it[len(it) - 1]
	t.Log(fmt.Sprintf("docker build -t awpaas/%s:%s .", name, version))
	t.Log(it)
}

//type Module struct {
//	Name, Group, Remarks, Version, Fid string
//	_id string
//}
//
//func (this *Module) check(args ...string) {
//	value := reflect.ValueOf(*this)
//	for _, arg := range args {
//		v := value.FieldByName(arg)
//		if !v.IsValid() {
//			break
//		}
//		if "" == v.Interface() {
//			fmt.Println("11111111111 null")
//			return
//		}
//	}
//}
//
//func Test_CheckModule(t *testing.T) {
//	module := &Module {
//		Name: "name",
//		Group: "group",
//		Remarks: "remarks",
//		Version : "version",
//		Fid: "fid",
//		_id: "12",
//	}
//	buf, err := json.Marshal(module)
//	if nil != err {
//		t.Error(err)
//	}
//	t.Log(string(buf))
//	module.check("Name", "Fid")
//}
