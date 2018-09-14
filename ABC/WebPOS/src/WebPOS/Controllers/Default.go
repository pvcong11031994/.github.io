package Controllers

import (
	"WebPOS/WebApp"
	"github.com/goframework/gf"
)

const (
	DEFAULT_FILTER_PATH string = "*"
)

func DefaultFilter(ctx *gf.Context) {
	ctx.ViewBases = []string{
		"master/master.html",
	}
	ctx.ViewData["WEBSITE_TITLE"] = ctx.Config.StrOrEmpty(WebApp.CONFIG_KEY_APP_WEBSITE_TITLE)
}

func Init() {
	gf.Filter(DEFAULT_FILTER_PATH, DefaultFilter)
	gf.Set404View("master/404.html")
}
