package Notice

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"unicode/utf8"
)

func ShowNoticeUpdate(ctx *gf.Context) {

	user := WebApp.GetContextUser(ctx)

	userNoticeModel := Models.UserNoticeModel{ctx.DB}

	notifyContent, err := userNoticeModel.GetNotice(user.ShopChainCd)
	Common.LogErr(err)
	ctx.ViewData["notify_content"] = notifyContent

	ctx.View = "maintenance/notice/update.html"
}

func UpdateNoticeUpdate(ctx *gf.Context) {
	user := WebApp.GetContextUser(ctx)

	userNoticeModel := Models.UserNoticeModel{ctx.DB}

	content := ctx.Form.StringNoEscape("content")

	if utf8.RuneCountInString(content) > 8000 {
		// FUJI-4677 EDIT
		//ctx.JsonResponse = map[string]interface{}{
		//	"is_success": false,
		//	"msg":        "文字が多すぎます。",
		//}
		ctx.JsonResponse = map[string]interface{}{
			"is_success": false,
			"msg":        "お知らせは8000文字以内で入力してください。",
		}
		return
	}

	err := userNoticeModel.SetNotice(user.ShopChainCd, user.UserID, content)
	if err != nil {
		Common.LogErr(err)

		// FUJI-4677 EDIT
		//ctx.JsonResponse = map[string]interface{}{
		//	"is_success": false,
		//	"msg":        "更新できません。",
		//}
		ctx.JsonResponse = map[string]interface{}{
			"is_success": false,
			"msg":        "更新に失敗しました。",
		}
		return
	}

	ctx.JsonResponse = map[string]interface{}{
		"is_success": true,
		"msg":        "お知らせが更新されました。",
	}
}

func ShowVJNoticeUpdate(ctx *gf.Context) {
	user := WebApp.GetContextUser(ctx)
	if user.ShopChainCd != Common.CHAIN_VJ {
		ctx.Redirect("/")
		return
	}

	userNoticeModel := Models.UserNoticeModel{ctx.DB}

	notifyContent, err := userNoticeModel.GetNotice("")
	Common.LogErr(err)
	ctx.ViewData["notify_content"] = notifyContent
	ctx.ViewData["vj_notify"] = true

	ctx.View = "maintenance/notice/update.html"
}

func UpdateVJNoticeUpdate(ctx *gf.Context) {
	user := WebApp.GetContextUser(ctx)
	if user.ShopChainCd != Common.CHAIN_VJ {
		ctx.Redirect("/")
		return
	}

	userNoticeModel := Models.UserNoticeModel{ctx.DB}
	content := ctx.Form.StringNoEscape("content")

	if utf8.RuneCountInString(content) > 8000 {
		// FUJI-4677 EDIT
		//ctx.JsonResponse = map[string]interface{}{
		//	"is_success": false,
		//	"msg":        "文字が多すぎます。",
		//}
		ctx.JsonResponse = map[string]interface{}{
			"is_success": false,
			"msg":        "お知らせは8000文字以内で入力してください。",
		}
		return
	}

	err := userNoticeModel.SetNotice("", user.UserID, content)
	if err != nil {
		Common.LogErr(err)

		// FUJI-4677 EDIT
		//ctx.JsonResponse = map[string]interface{}{
		//	"is_success": false,
		//	"msg":        "更新できません。",
		//}
		ctx.JsonResponse = map[string]interface{}{
			"is_success": false,
			"msg":        "更新に失敗しました。",
		}
		return
	}

	ctx.JsonResponse = map[string]interface{}{
		"is_success": true,
		"msg":        "お知らせが更新されました。",
	}
}
