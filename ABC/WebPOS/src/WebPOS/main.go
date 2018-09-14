package main

import (
	"WebPOS/Controllers"
	"WebPOS/Controllers/Account"
	"WebPOS/Controllers/Dashboard"
	"WebPOS/Controllers/Report"
	"WebPOS/ControllersApi"
	"github.com/goframework/gf"
	"log"
	"runtime"
	"WebPOS/Controllers/Maintenance"
	"WebPOS/Controllers/Download"
)

const (
	_SERVER_CONFIG = "../conf/server.cfg"
)

func main() {
	// Change config file
	gf.SetConfigPath(_SERVER_CONFIG)

	//
	log.Printf("runtime: %v %v %v", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	//Default filter & error handle
	Controllers.Init()

	// ログイン & Check role menu
	Account.Init()

	// ダッシュボード
	Dashboard.Init()

	// 帳票
	Report.Init()

	// Register handle route api
	Api.Init()

	// メンテナンスメニュー
	Maintenance.Init()

	// ダウンロードメニュー
	Download.Init()

	// Start job
	gf.Run()
}
