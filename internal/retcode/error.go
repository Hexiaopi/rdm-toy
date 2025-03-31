package retcode

import (
	"encoding/json"
	"fmt"
)

type RetErr struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

var codes = map[string]string{}

func NewCode(code, desc string) *RetErr {
	if v, ok := codes[code]; ok {
		panic(fmt.Sprintf("code:%s already exist with desc:%s", code, v))
	}
	codes[code] = desc
	return &RetErr{Code: code, Desc: desc}
}

func (e RetErr) Marshal() []byte {
	data, _ := json.Marshal(e)
	return data
}

func (e *RetErr) Error() string {
	return fmt.Sprintf("code:%s desc:%s", e.Code, e.Desc)
}
