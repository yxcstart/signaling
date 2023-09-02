package action

import (
	"fmt"
	"html/template"
	"net/http"
	"signaling/src/framework"

	log "github.com/sirupsen/logrus"
)

type xrtcClientPushAction struct{}

func NewXrtcClientPushAction() *xrtcClientPushAction {
	return &xrtcClientPushAction{}
}

func writeHtmlErrorResponse(w http.ResponseWriter, status int, errMsg string) {
	w.WriteHeader(status)
	w.Write([]byte(errMsg))
}

func (*xrtcClientPushAction) Execute(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(fmt.Sprintf("%s/template/push.html", framework.Conf.HttpStaticDir))
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
