package action

import (
	"encoding/json"
	"net/http"
	"signaling/src/common"
	"signaling/src/framework"
	"strconv"
)

const (
	CMDNO_PUSH      = 1
	CMDNO_PULL      = 2
	CMDNO_ANSWER    = 3
	CMDNO_STOP_PUSH = 4
	CMDNO_STOP_PULL = 5
)

type comHttpResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func errorResponse(cerr *common.Errors, w http.ResponseWriter, cr *framework.ComRequest) {
	cr.Logger.AddNotice("errCode", strconv.Itoa(cerr.ErrCode()))
	cr.Logger.AddNotice("errMsg", cerr.ErrorMsg())
	cr.Logger.Warningf("request process failed")

	rsp := comHttpResp{
		Code: cerr.ErrCode(),
		Msg:  "process error",
	}
	b, _ := json.Marshal(rsp)
	w.Write(b)
}
