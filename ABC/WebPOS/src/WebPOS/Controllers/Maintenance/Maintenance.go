package Maintenance

import (
	"WebPOS/Controllers/Maintenance/Users"
	"WebPOS/Controllers/Maintenance/Notice"
)

func Init() {
	// お知らせマスタメンテナンスメニュー
	Notice.Init()

	// ユーザマスタメンテナンスメニュー
	UsersMaintenance.Init()
}
