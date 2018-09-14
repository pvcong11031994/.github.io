package Account

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
)

const (
	LOGIN_ERROR_MESSAGE          = "ユーザーIDまたはパスワードが 正しくありません。"
	LOGIN_REJECT_MESSAGE         = "他の端末で利用中です。"
	PASSWORD_CHANGE_TIME_MESSAGE = "パスワード入力エラーが規定回数を超えています。パスワードの初期化をアカウント管理者に依頼してください。"
	SESSION_KEY_LOGIN_MSG        = "LOGIN_MSG"
	SESSION_KEY_LOGIN_ID         = "LOGIN_ID"
	LIMIT_STATUS_MESSAGE         = 30
	// ASO-5677 ログイン失敗回数を5から9に変更。(エンティティ定義がverchr 1のため。)
	// PASSWORD_CHANGE_TIME_LIMIT   = 5
	PASSWORD_CHANGE_TIME_LIMIT   = 9
)

type FormLogin struct {
	Id  string `form:"user"`
	Pwd string `form:"pwd"`
}

func Login(ctx *gf.Context) {
	userID := ctx.Session.Flashes(SESSION_KEY_LOGIN_ID)
	if len(userID) > 0 {
		ctx.ViewData["userID"] = userID[0]
		ctx.ViewData["focusPW"] = true
	} else {
		ctx.ViewData["focusID"] = true
	}
	msg := ctx.Session.Flashes(SESSION_KEY_LOGIN_MSG)
	if len(msg) > 0 {
		ctx.ViewData["msg"] = msg[0]
		if msg[0] != PASSWORD_CHANGE_TIME_MESSAGE {
			ctx.ViewData["flag"] = true
		} else {
			ctx.ViewData["flag"] = false
		}
	}

	// Get system status
	ssm := Models.SystemStatusModel{ctx.DB}
	listStatus := ssm.GetStatus(LIMIT_STATUS_MESSAGE)
	statusContent := ""
	for _, v := range listStatus {
		statusContent = statusContent + v.CreatedAt + "　" + v.Detail + "\r\n"
	}
	ctx.ViewData["status_content"] = statusContent

	// Get notice content
	userNoticeModel := Models.UserNoticeModel{ctx.DB}
	notifyContent, err := userNoticeModel.GetNotice("")
	Common.LogErr(err)
	ctx.ViewData["notify_content"] = Common.RenderLink("\n\n" + notifyContent)

	ctx.View = "account/login.html"
}

func LoginPost(ctx *gf.Context) {
	form := FormLogin{}
	ctx.Form.ReadStruct(&form)

	um := Models.UserMasterModel{ctx.DB}
	// check change pass time
	changePwTime, err := um.GetChangePWTime(form.Id)
	if changePwTime >= PASSWORD_CHANGE_TIME_LIMIT {
		ctx.Session.AddFlash(PASSWORD_CHANGE_TIME_MESSAGE, SESSION_KEY_LOGIN_MSG)
		ctx.Session.AddFlash(form.Id, SESSION_KEY_LOGIN_ID)
		ctx.Redirect(PATH_LOGIN)
		return
	}

	// check validate user, pw
	checkLogin, err := um.Login(form.Id, form.Pwd)
	Common.LogErr(err)

	if err != nil {
		Common.LogOutput(`"` +  form.Id +  `"` + " login fail.")
		ctx.Redirect(PATH_LOGIN)
	} else {
		if checkLogin {
			ctx.NewSession()
			user, err := um.GetUserInfoById(form.Id)
			Common.LogErr(err)

			if user != nil {
				ulsm := Models.UserLoginStatusModel{ctx.DB}
				err = ulsm.RejectLoggedInUser(user.UserID)
				Common.LogErr(err)
				ipAddress := ctx.GetRequestForwardedIP()
				if ipAddress == "" {
					ipAddress = ctx.GetRequestIP()
				}
				browserInfo := ctx.GetBrowserAgent()
				deviceID := "" //TODO: login device control processing
				err = ulsm.NewLoginStatus(user.UserID, ctx.Session.ID, ipAddress, browserInfo, deviceID)
				Common.LogErr(err)

				rpm := Models.ReportDesignMaster{ctx.DB}
				if design, err := rpm.GetReportDesign(user.ShopChainCd); err != nil {
					Common.LogErr(err)
				} else {
					user.Design = design
				}
				WebApp.SetSessionUser(ctx, user)
			}
			// update um_pw_change_flg = 0
			err = um.UpdateChangePWFlag(form.Id, false)
			Common.LogErr(err)
			Common.LogOutput(`"` +  user.UserName +  `"` + " login success.")
			ctx.Redirect(PATH_ROOT)
		} else {
			// update um_pw_change_flg Auto Increment when input error pass
			err := um.UpdateChangePWFlag(form.Id, true)
			Common.LogErr(err)
			if changePwTime+1 >= PASSWORD_CHANGE_TIME_LIMIT {
				ctx.Session.AddFlash(PASSWORD_CHANGE_TIME_MESSAGE, SESSION_KEY_LOGIN_MSG)
			} else {
				ctx.Session.AddFlash(LOGIN_ERROR_MESSAGE, SESSION_KEY_LOGIN_MSG)

			}
			ctx.Session.AddFlash(form.Id, SESSION_KEY_LOGIN_ID)
			Common.LogOutput(`"` +  form.Id +  `"` + " login fail.")
			ctx.Redirect(PATH_LOGIN)
		}
	}
}

func Logout(ctx *gf.Context) {
	user := WebApp.GetContextUser(ctx)
	ulsm := Models.UserLoginStatusModel{ctx.DB}
	ulsm.LogoutByKey(user.UserID, ctx.Session.ID)

	ctx.ClearSession()
	Common.LogOutput(`"` +  user.UserName +  `"` + " logout.")
	ctx.Redirect(PATH_LOGIN)
}
