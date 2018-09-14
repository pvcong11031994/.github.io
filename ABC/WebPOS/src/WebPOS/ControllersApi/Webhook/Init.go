package Webhook

import (
	"WebPOS/ControllersApi/Utils"
	"github.com/goframework/gf"
)

const (
	ROUTE_API_JOB = ApiUtils.ROUTE_API + "/webhook"
)

//webhookAPIを初期化
func Init()  {

	gf.HandlePost(ROUTE_API_JOB, doProcess)
}