package Dashboard

import (
	"WebPOS/Common"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Redirect to top page ("/top")
func Dashboard(ctx *gf.Context) {

	ctx.Redirect(PATH_TOP_PAGE)
}

// Show top page html
func TopPage(ctx *gf.Context) {

	ctx.View = "dashboard/top_page.html"

	// Load path of index.html
	dir := ctx.Config.Str(WebApp.CONFIG_KEY_APP_STATIC_DIR, _DEFAULT_FOLDER_TOP_PAGE)
	pathDir, err := filepath.Abs(dir)
	if err != nil {
		log.Println(err)
		return
	}
	filepathHtml := filepath.Join(pathDir, _DEFAULT_TOP_PAGE)

	data, err := ioutil.ReadFile(filepathHtml)
	if err != nil {
		log.Println(err)
		return
	} else if data == nil {
		return
	}
	htmlData := string(data)

	// Read all file in dir
	listFile, err := ioutil.ReadDir(pathDir)
	Common.LogErr(err)

	desSrc, err := filepath.Abs(ctx.Config.Str(gf.CFG_KEY_SERVER_STATIC_DIR, _DEFAULT_STATIC_DIR))
	Common.LogErr(err)
	dirTopPage := filepath.Join(desSrc, _DEFAULT_FOLDER_STATIC_TOP_PAGE)
	err = os.RemoveAll(dirTopPage)
	Common.LogErr(err)
	err = os.MkdirAll(dirTopPage, os.ModePerm)
	Common.LogErr(err)

	// copy file img, css, js,... into SERVER_STATIC_DIR
	for _, file := range listFile {
		if !strings.HasSuffix(file.Name(), _SUNFIX_HTML) &&
			!strings.HasSuffix(file.Name(), _SUNFIX_HTM) &&
			strings.Contains(htmlData, file.Name()) {
			htmlData = strings.Replace(htmlData,
				file.Name(),
				filepath.Join(_DEFAULT_FOLDER_STATIC, _DEFAULT_FOLDER_STATIC_TOP_PAGE, file.Name()),
				-1,
			)
			src := filepath.Join(pathDir, file.Name())
			des := filepath.Join(dirTopPage, file.Name())
			copyFile(src, des)
		}
	}

	ctx.ViewData["index"] = htmlData
}

func copyFile(src string, dst string) {

	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	Common.LogErr(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	Common.LogErr(err)
}
