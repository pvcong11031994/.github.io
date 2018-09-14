package Api

import (
	"WebPOS/ControllersApi/Init"
	"WebPOS/ControllersApi/User"
	"WebPOS/ControllersApi/Webhook"
)

// Register handle all route controller
func Init() {
	//
	ApiInit.FilterIp()
	//
	User.Init()
	//
	Webhook.Init()
}
