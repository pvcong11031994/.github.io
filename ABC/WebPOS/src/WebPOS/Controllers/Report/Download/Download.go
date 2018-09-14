package Download

import (
	"github.com/goframework/gf"
)

func Download(ctx *gf.Context) {

	f := ctx.Form.String("f")
	fileCSV := f
	ctx.ServeStaticFile(fileCSV, true)
}
