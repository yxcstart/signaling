package framework

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	http.HandleFunc("/", entry)
}

type ActionInterface interface {
	Execute(w http.ResponseWriter, r *http.Request)
}

type ComRequest struct {
	R      *http.Request
	Logger *ComLog
	LogId  uint32
}

var GActionRouter map[string]ActionInterface = map[string]ActionInterface{}

func responseError(w http.ResponseWriter, r *http.Request, status int, errMsg string) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf("%d - %s", status, errMsg)))
}

func getRealClientIP(r *http.Request) string {
	ip := r.RemoteAddr
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		ip = rip
	} else if rip = r.Header.Get("X-Forwarded-IP"); rip != "" {
		ip = rip
	}

	return ip
}

func entry(w http.ResponseWriter, r *http.Request) {
	if "/favicon.ico" == r.URL.Path {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
		return
	}

	if action, ok := GActionRouter[r.URL.Path]; ok {
		if action != nil {
			cr := &ComRequest{
				R:      r,
				Logger: &ComLog{},
				LogId:  GetLogId32(),
			}

			start := time.Now()
			cr.Logger.AddNotice("logId", strconv.Itoa(int(cr.LogId)))
			cr.Logger.AddNotice("url", r.URL.Path)
			cr.Logger.AddNotice("referer", r.Header.Get("Referer"))
			cr.Logger.AddNotice("cookie", r.Header.Get("Cookie"))
			cr.Logger.AddNotice("ua", r.Header.Get("User-Agent"))
			cr.Logger.AddNotice("clientIP", r.RemoteAddr)
			cr.Logger.AddNotice("realClientIP", getRealClientIP(r))
			r.ParseForm()
			action.Execute(w, r)
			cr.Logger.AddNotice("cost", time.Since(start).String())
			cr.Logger.Infof("")
		} else {
			responseError(w, r, http.StatusInternalServerError, "internal server error")
		}
	} else {
		responseError(w, r, http.StatusNotFound, "not found")
	}

}

func StartHttp() error {
	log.Infof("start http...")
	return http.ListenAndServe(fmt.Sprintf(":%d", Conf.HttpPort), nil)
}
