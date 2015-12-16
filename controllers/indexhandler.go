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

func (self *IndexHandler) Handle(w http.ResponseWriter, r *http.Request, g kernel.G, data IndexedData) {
	searcher, _ := g.DIY["searcher"].(search.Engine)
	searcher.IndexDocument(data.DocId, search.DocumentIndexData{
		Content: data.Content,
		Labels:  data.Labels,
	})
	self.JsonResponse(w, "", 200)
}

func (self *IndexHandler) checkArgument(text string, docid string, labels string) (ok bool, data IndexedData) {
	if text == "" || docid == "" {
		return
	} else {
		data.Content = text
		id, err := strconv.Atoi(docid)
		if err != nil {
			return
		} else {
			data.DocId = uint64(id)
			ok = true
			return
		}
		data.Labels = strings.Split(labels, "-")
		ok = true
		return
	}
}

func (self *IndexHandler) Post(w http.ResponseWriter, r *http.Request, g kernel.G) {
	// 得到要分词的文本
	r.ParseForm()
	text := r.PostFormValue("text")
	docid := r.PostFormValue("docid")
	labels := r.PostFormValue("tags")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	log.Println(text, docid, labels)
	if ok, data := self.checkArgument(text, docid, labels); !ok {
		self.JsonResponse(w, nil, 401)
	} else {
		self.Handle(w, r, g, data)
	}
}

func (self *IndexHandler) Get(w http.ResponseWriter, r *http.Request, g kernel.G) {
	text := r.URL.Query().Get("text")
	docid := r.URL.Query().Get("docid")
	labels := r.URL.Query().Get("tags")
	log.Printf("Method: %s From Ip: %s", r.Method, r.RemoteAddr)
	if ok, data := self.checkArgument(text, docid, labels); !ok {
		self.JsonResponse(w, nil, 401)
	} else {
		self.Handle(w, r, g, data)
	}
}
