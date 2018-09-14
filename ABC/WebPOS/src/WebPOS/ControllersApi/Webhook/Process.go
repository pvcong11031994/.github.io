package Webhook

import (
	"github.com/goframework/gf"
	"net/http"
	"WebPOS/Models/DB"
)

func doProcess(ctx *gf.Context) {

	formInput := RequestDto{}
	ctx.Form.ReadStruct(&formInput)

	ssModel := Models.SystemStatusModel{ctx.DB}
	err := ssModel.InsertStatus(formInput.Chain, formInput.Group, formInput.Detail)

	if err != nil {
		ctx.JsonResponse = ResponseDto{
			Success: false,
			Code:    http.StatusInternalServerError,
			Msg:     err.Error(),
		}
	} else {
		ctx.JsonResponse = ResponseDto{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "success",
		}
	}
}