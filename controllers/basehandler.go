package controllers

import (
	"encoding/json"
	"github.com/aosen/kernel"
	"io"
	"net/http"
)

/*返回的JSON数据*/
type Response struct {
	Values interface{} `json:"result"`
	Err    int         `json:"err"`
	Errmsg string      `json:"errmsg"`
}

var ERR = map[int]string{
	0: "success",
	1: "Invalid argument",
}

type BaseHandler struct {
	kernel.W
}

func (self *BaseHandler) JsonResponse(w http.ResponseWriter, v interface{}, code int) {
	resp, _ := json.Marshal(&Response{
		Values: v,
		Err:    code,
		Errmsg: ERR[code],
	})
	//w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(resp))
}
