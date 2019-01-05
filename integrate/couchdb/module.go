package couchdb

import (
	"../../util"
)

type condition map[string]interface{}

func (this *condition) String() string {
	reply, _ := util.FormatToString(this)
	return reply
}

func Condition() *condition {
	return &condition{"selector": condition{}, "fields": "", "limit": 10, "skip": 0}
}

func (this *condition) Append(key string, args ...interface{}) *condition {
	selector := (*this)["selector"].(condition)
	c := condition{}
	for i := 0; i < len(args) - 1; i+=2 {
		c[args[i].(string)] = args[i+1]
	}
	selector[key] = c
	return this
}

func (this *condition) Fields(args ...string) *condition {
	(*this)["fields"] = args
	return this
}

func (this *condition) Page(begin, limit int) *condition {
	(*this)["skip"] = begin
	(*this)["limit"] = limit
	return this
}
