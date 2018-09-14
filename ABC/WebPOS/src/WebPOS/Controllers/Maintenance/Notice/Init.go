package Notice

import (
	"github.com/goframework/gf"
)

const (
	PATH_MAINTENANCE_NOTICE_UPDATE    string = "/maintenance/notice/update"
	PATH_MAINTENANCE_VJ_NOTICE_UPDATE string = "/maintenance/notice/vj/update"
)

func Init() {
	// お知らせ
	gf.HandleGet(PATH_MAINTENANCE_NOTICE_UPDATE, ShowNoticeUpdate)
	gf.HandlePost(PATH_MAINTENANCE_NOTICE_UPDATE, UpdateNoticeUpdate)
	// VJお知らせ
	gf.HandleGet(PATH_MAINTENANCE_VJ_NOTICE_UPDATE, ShowVJNoticeUpdate)
	gf.HandlePost(PATH_MAINTENANCE_VJ_NOTICE_UPDATE, UpdateVJNoticeUpdate)
}
