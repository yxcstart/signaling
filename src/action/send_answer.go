package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"signaling/src/common"
	"signaling/src/framework"
	"strconv"
)

type sendAnswerAction struct{}

func NewSendAnswerAction() *sendAnswerAction {
	return &sendAnswerAction{}
}

type sendAnswerReq struct {
	Cmdno      int    `json:"cmdno"`
	Uid        uint64 `json:"uid"`
	StreamName string `json:"stream_name"`
	Answer     string `json:"answer"`
	Type       string `json:"type"`
}

type sendAnswerResp struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
}

func (*sendAnswerAction) Execute(w http.ResponseWriter, cr *framework.ComRequest) {
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

	// answer
	var answer string
	if values, ok := r.Form["answer"]; ok {
		answer = values[0]
	}

	if "" == answer {
		cerr := common.New(common.ParamErr, "answer is null")
		errorResponse(cerr, w, cr)
		return
	}

	// type
	var strType string
	if values, ok := r.Form["type"]; ok {
		strType = values[0]
	}

	if "" == strType {
		cerr := common.New(common.ParamErr, "strType is null")
		errorResponse(cerr, w, cr)
		return
	}

	req := sendAnswerReq{
		Cmdno:      CMDNO_ANSWER,
		Uid:        uid,
		StreamName: streamName,
		Answer:     answer,
		Type:       strType,
	}

	var resp sendAnswerResp

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
	}

	b, _ := json.Marshal(httpResp)
	cr.Logger.AddNotice("resp", string(b))
	w.Write(b)
}
