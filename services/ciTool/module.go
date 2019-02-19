package ciTool

import "fmt"

var registryType = [4]string{}

type ci struct {
	Id int64
	WarehouseId int64 `json:"warehouseId"`
	RegistryType string `json:"registryType"`
	Content []*cmd `json:"content"`
	LastReport string `json:"lastReport"`
	LastCiTime int64 `json:"lastCiTime"`
	CreateTime int64 `json:"createTime"`
	Status bool
}

type cmd struct {
	Id int64
	Context string
	CreateTime int64
	ModifyTime int64
	Status bool
	CId int64
}

func (this *ci) String() string {
	return fmt.Sprintf("{'fileType': '%s', 'content': '%v'}", this.RegistryType, this.Content)
}

func (this *cmd) String() string {
	return fmt.Sprintf("{'context': '%s', 'createTime': '%s'}", this.Context, this.CreateTime)
}