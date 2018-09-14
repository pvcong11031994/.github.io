package Dashboard

import (
	"github.com/goframework/gf"
)

const (
	PATH_ROOT     = "/"
	PATH_TOP_PAGE = "/top"

	_DEFAULT_FOLDER_TOP_PAGE        = "/var/ba/toppage/"
	_DEFAULT_TOP_PAGE               = "index.html"
	_DEFAULT_STATIC_DIR             = "./static"
	_DEFAULT_FOLDER_STATIC_TOP_PAGE = "top_page"
	_DEFAULT_FOLDER_STATIC          = "static"
	_SUNFIX_HTML                    = "html"
	_SUNFIX_HTM                     = "htm"
)

func Init() {

	gf.HandleGetPost(PATH_ROOT, Dashboard)
	gf.HandleGetPost(PATH_TOP_PAGE, TopPage)
}
