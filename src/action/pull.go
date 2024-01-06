package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"signaling/src/common"
	"signaling/src/framework"
	"strconv"
)

type pullAction struct{}

func NewPullAction() *pullAction {
	return &pullAction{}
}

type xrtcPullReq struct {
	Cmdno      int    `json:"cmdno"`
	Uid        uint64 `json:"uid"`
	StreamName string `json:"stream_name"`
	Audio      int    `json:"audio"`
	Video      int    `json:"video"`
}

type xrtcPullResp struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Offer  string `json:"offer"`
}

type pullData struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

func (*pullAction) Execute(w http.ResponseWriter, cr *framework.ComRequest) {
	r := cr.R
	var strUid string
	if val, ok := r.Form["uid"]; ok {
		strUid = val[0]
	}

	uid, err := strconv.ParseUint(strUid, 10, 64)
	if err != nil || uid <= 0 {
		cerr := common.New(common.ParamErr, "parse uid err : "+err.Error())
		errorResponse(cerr, w, cr)
		return
	}

	// streamName
	var streamName string
	if values, ok := r.Form["streamName"]; ok {
		streamName = values[0]
	}

	if "" == streamName {
		cerr := common.New(common.ParamErr, "streamName is null")
		errorResponse(cerr, w, cr)
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

	req := xrtcPullReq{
		Cmdno:      CMDNO_PULL,
		Uid:        uid,
		StreamName: streamName,
		Audio:      audio,
		Video:      video,
	}

	var resp xrtcPullResp

	err = framework.Call("xrtc", req, &resp, cr.LogId)
	if err != nil {
		cerr := common.New(common.NetworkErr, "backend process error"+err.Error())
		errorResponse(cerr, w, cr)
		return
	}

	if resp.ErrNo != 0 {
		cerr := common.New(common.NetworkErr, fmt.Sprintf("backend process errno: %d", resp.ErrNo))
		errorResponse(cerr, w, cr)
		return
	}

	httpResp := comHttpResp{
		Code: 0,
		Msg:  "success",
		Data: pushData{
			Type: "offer",
			Sdp:  resp.Offer,
		},
	}

	b, _ := json.Marshal(httpResp)
	cr.Logger.AddNotice("resp", string(b))
	w.Write(b)
}
