package controllers

import (
	"fmt"
	"github.com/aosen/kernel"
	"log"
	"net/http"
)

type NotFoundHandler struct {
	BaseHandler
}

func (self *NotFoundHandler) Prepare(w http.ResponseWriter, r *http.Request, g kernel.G) {
	log.Printf("%s 404 not found %s", r.Method, r.URL)
	fmt.Fprintln(w, "sorry 404 not found")
}
