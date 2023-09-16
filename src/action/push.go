package action

import (
	"net/http"
	"signaling/src/common"
	"signaling/src/framework"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type pushAction struct{}

func NewPushAction() *pushAction {
	return &pushAction{}
}

func (*pushAction) Execute(w http.ResponseWriter, cr *framework.ComRequest) {
	r := cr.R
	var strUid string
	if val, ok := r.Form["uid"]; ok {
		strUid = val[0]
	}

	uid, err := strconv.ParseUint(strUid, 10, 64)
	if err != nil || uid <= 0 {
		cerr := common.New(common.ParamErr, "parse uid err : "+err.Error())
		common.ErrorResponse(cerr, w, cr)
		return
	}

	// streamName
	var streamName string
	if values, ok := r.Form["streamName"]; ok {
		streamName = values[0]
	}

	if "" == streamName {
		cerr := common.New(common.ParamErr, "streamName is null")
		common.ErrorResponse(cerr, w, cr)
		return
	}

	// audio video
	var strAudio, strVideo string
	var audio, video int

	if values, ok := r.Form["audio"]; ok {
		strAudio = values[0]
	}

	if "" == strAudio || "0" == strAudio {
		audio = 0
	} else {
		audio = 1
	}

	if values, ok := r.Form["video"]; ok {
		strVideo = values[0]
	}

	if "" == strVideo || "0" == strVideo {
		video = 0
	} else {
		video = 1
	}

	log.Infof("push param uid:%d, streamName:%s, video:%d, audio:%d", uid, streamName, video, audio)
}
