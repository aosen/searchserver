/*
cut 分词服务器提供JSON格式的RPC服务
datatime: 2015-10-12
author: aosen
homepage: https://github.com/aosen/
url:
    "/"
输入:
    GET模式输入text参数
输出:
    {
        result:[
            {"text":"服务器", "pos":"n"},
            {"text":"指令", "pos":"n"},
            ...
        ],
        err: 0,
        errmsg: ""
    }
*/
package controllers

import (
	"github.com/aosen/cut"
	"github.com/aosen/kernel"
	"log"
	"net/http"
	"strings"
)

type CutHandler struct {
	BaseHandler
}

var (
	//生成分词器对象
	segmenter = cut.Segmenter{}
)

func (self *CutHandler) Base(w http.ResponseWriter, r *http.Request, g kernel.G, text string, mode string) {
	type Value struct {
		Text string `json:"text"`
		Pos  string `json:"pos"`
	}
	segmenter, _ := g.DIY["seg"].(cut.Segmenter)
	//通用处理方法
	if text == "" || (!strings.EqualFold(mode, "1") && !strings.EqualFold(mode, "0")) {
		self.JsonResponse(w, nil, 401)
	} else {
		//整理为输出格式
		s := []*Value{}
		//开始分词
		func() {
			for _, seg := range segmenter.Cut([]byte(text),
				func() bool {
					if mode == "1" {
						return true
					} else {
						return false
					}
				}()) {
				s = append(s, &Value{
					Text: seg.Token().Text(),
					Pos:  seg.Token().Pos(),
				})
			}
		}()
		self.JsonResponse(w, s, 200)
	}
}

func (self *CutHandler) Post(w http.ResponseWriter, r *http.Request, g kernel.G) {
	// 得到要分词的文本
	r.ParseForm()
	text := r.PostFormValue("text")
	mode := r.PostFormValue("mode")
	log.Printf("Method: %s From Ip: %s text: %s mode: %s", r.Method, r.RemoteAddr, text, mode)
	self.Base(w, r, g, text, mode)
}

func (self *CutHandler) Get(w http.ResponseWriter, r *http.Request, g kernel.G) {
	text := r.URL.Query().Get("text")
	mode := r.URL.Query().Get("mode")
	log.Printf("Method: %s From Ip: %s text: %s mode: %s", r.Method, r.RemoteAddr, text, mode)
	self.Base(w, r, g, text, mode)
}
