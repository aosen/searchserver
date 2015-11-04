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

var sr SearchRequest

func (self *SearchHandler) Handle(w http.ResponseWriter, r *http.Request, g kernel.G) {
	searcher, _ := g.DIY["searcher"].(search.Engine)
	self.JsonResponse(w, searcher.Search(search.SearchRequest{
		Text:    sr.Text,
		Labels:  sr.Labels,
		DocIds:  sr.DocIds,
		Timeout: sr.TimeOut,
	}), 0)
}

func (self *SearchHandler) checkArgument(text string, docids string, labels string, timeout string) bool {
	//检测参数合法性，合法返回True 不合法返回False
	if text == "" && docids == "" && labels == "" {
		return false
	} else {
		if timeout != "" {
			t, err := strconv.Atoi(timeout)
			if err != nil {
				return false
			}
			sr.TimeOut = t
		}
		if text != "" {
			sr.Text = text
		}
		if labels != "" {
			sr.Labels = strings.Split(labels, "-")
		}
		if docids != "" {
			for _, id := range strings.Split(docids, "-") {
				tmp, err := strconv.Atoi(id)
				if err != nil {
					return false
				} else {
					sr.DocIds = append(sr.DocIds, uint64(tmp))
				}
			}
		}
	}
	log.Println("test")
	return true
}

func (self *SearchHandler) Post(w http.ResponseWriter, r *http.Request, g kernel.G) {
	// 得到要分词的文本
	r.ParseForm()
	text := r.PostFormValue("text")
	docids := r.PostFormValue("docids")
	labels := r.PostFormValue("tags")
	timeout := r.PostFormValue("timeout")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	if !self.checkArgument(text, docids, labels, timeout) {
		self.JsonResponse(w, nil, 1)
	} else {
		self.Handle(w, r, g)
	}
}

func (self *SearchHandler) Get(w http.ResponseWriter, r *http.Request, g kernel.G) {
	text := r.URL.Query().Get("text")
	docids := r.URL.Query().Get("docids")
	labels := r.URL.Query().Get("tags")
	timeout := r.URL.Query().Get("timeout")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	if !self.checkArgument(text, docids, labels, timeout) {
		self.JsonResponse(w, nil, 1)
	} else {
		self.Handle(w, r, g)
	}
}
