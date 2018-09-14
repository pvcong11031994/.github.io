package Download

import (
	"github.com/goframework/gf"
)

const (
	PATH_REPORT_DOWN_LOAD string = "/report/download"
)

func Init() {
	// ダウンロード
	gf.HandlePost(PATH_REPORT_DOWN_LOAD, Download)
}
