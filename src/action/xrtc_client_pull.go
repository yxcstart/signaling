package action

import (
	"fmt"
	"html/template"
	"net/http"
	"signaling/src/framework"

	log "github.com/sirupsen/logrus"
)

type xrtcClientPullAction struct{}

func NewXrtcClientPullAction() *xrtcClientPullAction {
	return &xrtcClientPullAction{}
}

func (*xrtcClientPullAction) Execute(w http.ResponseWriter, cr *framework.ComRequest) {
	r := cr.R
	t, err := template.ParseFiles(fmt.Sprintf("%s/template/pull.html", framework.Conf.HttpStaticDir))
	if err != nil {
		log.Errorf(err.Error())
		writeHtmlErrorResponse(w, http.StatusNotFound, "404 - not found")
		return
	}

	request := make(map[string]string)

	for k, v := range r.Form {
		request[k] = v[0]
	}

	err = t.Execute(w, request)
	if err != nil {
		log.Errorf(err.Error())
		writeHtmlErrorResponse(w, http.StatusNotFound, "404 - not found")
		return
	}
}
