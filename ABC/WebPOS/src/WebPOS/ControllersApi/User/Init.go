package User

import (
	"WebPOS/ControllersApi/Utils"
	"github.com/goframework/gf"
)

const (
	ROUTE_API_USER = ApiUtils.ROUTE_API + "/user"
)

// Init route route api
func Init() {
	gf.HandleGetPost(ROUTE_API_USER, doProcess)
}
