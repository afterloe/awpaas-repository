package ciTool

import "fmt"

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