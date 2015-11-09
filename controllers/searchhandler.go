package controllers

import (
	"github.com/aosen/kernel"
	"github.com/aosen/search"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SearchHandler struct {
	BaseHandler
}

type SearchRequest struct {
	DocIds  []uint64
	Text    string
	Labels  []string
	TimeOut int
}

func (self *SearchHandler) Handle(w http.ResponseWriter, r *http.Request, g kernel.G, sr SearchRequest) {
	searcher, _ := g.DIY["searcher"].(search.Engine)
	self.JsonResponse(w, searcher.Search(search.SearchRequest{
		Text:    sr.Text,
		Labels:  sr.Labels,
		DocIds:  sr.DocIds,
		Timeout: sr.TimeOut,
	}), 200)
}

func (self *SearchHandler) checkArgument(text string, docids string, labels string, timeout string) (ok bool, sr SearchRequest) {
	//检测参数合法性，合法返回True 不合法返回False
	if text == "" && labels == "" || docids == "" {
		return
	} else {
		if timeout != "" {
			t, err := strconv.Atoi(timeout)
			if err != nil {
				return
			}
			sr.TimeOut = t
		}
		if text != "" {
			sr.Text = text
		}
		if labels != "" {
			sr.Labels = strings.Split(labels, "-")
		}
		//生成docids范围内的数组切片
		docidlist := strings.Split(docids, "-")
		if len(docidlist) != 2 {
			return
		} else {
			for _, id := range strings.Split(docids, "-") {
				tmp, err := strconv.Atoi(id)
				if err != nil {
					return
				} else {
					sr.DocIds = append(sr.DocIds, uint64(tmp))
				}
			}
			log.Println(sr.DocIds)
		}
	}
	ok = true
	return
}

func (self *SearchHandler) Post(w http.ResponseWriter, r *http.Request, g kernel.G) {
	// 得到要分词的文本
	r.ParseForm()
	text := r.PostFormValue("text")
	docids := r.PostFormValue("docids")
	labels := r.PostFormValue("tags")
	timeout := r.PostFormValue("timeout")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	if ok, sr := self.checkArgument(text, docids, labels, timeout); !ok {
		self.JsonResponse(w, nil, 401)
	} else {
		self.Handle(w, r, g, sr)
	}
}

func (self *SearchHandler) Get(w http.ResponseWriter, r *http.Request, g kernel.G) {
	text := r.URL.Query().Get("text")
	docids := r.URL.Query().Get("docids")
	labels := r.URL.Query().Get("tags")
	timeout := r.URL.Query().Get("timeout")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	if ok, sr := self.checkArgument(text, docids, labels, timeout); !ok {
		self.JsonResponse(w, nil, 401)
	} else {
		self.Handle(w, r, g, sr)
	}
}
