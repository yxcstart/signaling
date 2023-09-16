package action

import (
	"net/http"
	"signaling/src/framework"
)

type pushAction struct{}

func NewPushAction() *pushAction {
	return &pushAction{}
}

func (*pushAction) Execute(w http.ResponseWriter, cr *framework.ComRequest) {

}
