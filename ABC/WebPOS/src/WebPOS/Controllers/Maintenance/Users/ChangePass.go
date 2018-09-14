package UsersMaintenance

import (
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
)

func ChangePassView(ctx *gf.Context) {
	ctx.ViewData["link_action"] = PATH_MAINTENANCE_USER_CHANGE_PASS_ACTION
	ctx.View = "maintenance/user/change_pass.html"
}

func ChangePassAction(ctx *gf.Context) {
	data := map[string]interface{}{}
	data["is_success"] = "false"

	newPass := ctx.Form.String("pwd_new")
	renewPass := ctx.Form.String("pwd_renew")
	if newPass != renewPass {
		// FUJI-4677 EDIT
		//data["message_err"] = "「新規パスワード」と「確認入力」を同じ入力してください。"
		data["message_err"] = "「新規パスワード」と「確認入力」は同じ内容を入力してください。"
		ctx.JsonResponse = data
		return
	}

	user := WebApp.GetContextUser(ctx)

	um := Models.UserMasterModel{ctx.DB}
	cfPwExpiringDate := ctx.Config.Int(WebApp.CONFIG_KEY_PW_EXPIRING_DATE, 30)
	if err := um.ChangePass(user.UserID, newPass, cfPwExpiringDate); err != nil {
		data["message_err"] = "Error : " + err.Error()
		data["is_success"] = "false"
	} else {
		newUseInfo, _ := um.GetUserInfoById(user.UserID)
		WebApp.SetSessionUser(ctx, newUseInfo)
		data["is_success"] = "true"
	}

	ctx.JsonResponse = data
}
