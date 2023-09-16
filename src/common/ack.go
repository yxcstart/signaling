package common

import (
	"encoding/json"
	"net/http"
	"signaling/src/framework"
	"strconv"
)

type comHttpResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ErrorResponse(cerr *Errors, w http.ResponseWriter, cr *framework.ComRequest) {
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
