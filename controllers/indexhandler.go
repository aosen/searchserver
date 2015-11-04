package controllers

import (
	"github.com/aosen/kernel"
	"github.com/aosen/search"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type IndexedData struct {
	DocId   uint64
	Content string
	Labels  []string
}

type IndexHandler struct {
	BaseHandler
}

var data IndexedData

func (self *IndexHandler) Handle(w http.ResponseWriter, r *http.Request, g kernel.G) {
	searcher, _ := g.DIY["searcher"].(search.Engine)
	searcher.IndexDocument(data.DocId, search.DocumentIndexData{
		Content: data.Content,
		Labels:  data.Labels,
	})
	self.JsonResponse(w, "", 0)
}

func (self *IndexHandler) checkArgument(text string, docid string, labels string) bool {
	if text == "" || docid == "" {
		return false
	} else {
		data.Content = text
		id, err := strconv.Atoi(docid)
		if err != nil {
			return false
		} else {
			data.DocId = uint64(id)
			return true
		}
		data.Labels = strings.Split(labels, "-")
		return true
	}
}

func (self *IndexHandler) Post(w http.ResponseWriter, r *http.Request, g kernel.G) {
	// 得到要分词的文本
	r.ParseForm()
	text := r.PostFormValue("text")
	docid := r.PostFormValue("docid")
	labels := r.PostFormValue("tags")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	if !self.checkArgument(text, docid, labels) {
		self.JsonResponse(w, nil, 1)
	} else {
		self.Handle(w, r, g)
	}
}

func (self *IndexHandler) Get(w http.ResponseWriter, r *http.Request, g kernel.G) {
	text := r.URL.Query().Get("text")
	docid := r.URL.Query().Get("docid")
	labels := r.URL.Query().Get("tags")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	if !self.checkArgument(text, docid, labels) {
		self.JsonResponse(w, nil, 1)
	} else {
		self.Handle(w, r, g)
	}
}
