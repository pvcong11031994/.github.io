package UsersMaintenance

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"fmt"
	"github.com/goframework/gf"
	"strings"
)

func New(ctx *gf.Context) {

	user := WebApp.GetContextUser(ctx)

	sm := Models.ShopMasterModel{ctx.DB}
	um := Models.UserMasterModel{ctx.DB}
	cm := Models.CorpMasterModel{ctx.DB}

	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	listMenu := um.ListMenuGroup(user.ShopChainCd, user.ShopCd)

	sfgm := Models.ShopFranchiseGroupMasterModel{ctx.DB}
	listFranchiseCd := []ModelItems.FranchiseCdItem{}
	err = sfgm.GetListFranchiseCD(user.FranchiseGroupCd, &listFranchiseCd)
	Common.LogErr(err)

	listCorp := []ModelItems.CorpItem{}
	err = cm.GetListCorp(user.CorpCd, &listCorp)
	Common.LogErr(err)

	ctx.ViewData["corps"] = listCorp
	ctx.ViewData["shops"] = listShop
	ctx.ViewData["menugroups"] = listMenu
	ctx.ViewData["listFranchiseCd"] = listFranchiseCd
	ctx.ViewData["link_insert_check"] = PATH_MAINTENANCE_USER_NEW_CHECK
	ctx.ViewData["link_result_confirm"] = PATH_MAINTENANCE_USER_CONFIRM
	ctx.View = "maintenance/user/new.html"
}

func CheckAccount(ctx *gf.Context) {
	//▼▼▼▼▼▼▼▼▼
	data := map[string]interface{}{}
	list_user := ctx.Form.String("user_id")
	list_user = strings.TrimRight(list_user, ",")
	l_user := strings.Split(list_user, ",")
	//▲▲▲▲▲▲▲▲▲
	//▼▼▼▼▼▼▼▼▼
	um := Models.UserMasterModel{ctx.DB}
	var list_exist_id []int
	data["is_success"] = true

	for i, val := range l_user {
		exist_id, err := um.CheckUser(val)
		Common.LogErr(err)
		if exist_id && val != "" {
			list_exist_id = append(list_exist_id, i)
			data["is_success"] = false
		}
	}
	//▲▲▲▲▲▲▲▲▲
	//▼▼▼▼▼▼▼▼▼

	data["list_id"] = list_exist_id
	ctx.JsonResponse = data
	//▲▲▲▲▲▲▲▲▲
}

func Confirm(ctx *gf.Context) {
	//▼▼▼▼▼▼▼▼▼
	//パラメーターから取得する
	//ユーザIDから取得する
	user_id := ctx.Form.Array("user_id")
	//ユーザ名から取得する
	user_name := ctx.Form.Array("user_name")
	//店舗から取得する
	shop_cd := ctx.Form.Array("shop_cd")
	//権限から取得する
	flg_auth := ctx.Form.Array("flg_auth")
	//フランチャイズコードから取得する
	franchise_cd := ctx.Form.Array("franchise_cd")
	//メニューグループから取得する
	flg_menu_group := ctx.Form.Array("flg_menu_group")
	//企業コードから取得する
	corp_cd := ctx.Form.Array("corp")
	//部署コードから取得する
	dept_cd := ctx.Form.Array("dept_cd")
	//所属部署名から取得する
	dept_name := ctx.Form.Array("dept_name")
	//メールから取得する
	user_mail := ctx.Form.Array("user_mail")
	//電話から取得する
	user_phone := ctx.Form.Array("user_phone")
	//FAXから取得する
	user_xerox := ctx.Form.Array("user_xerox")
	//パスワードから取得する
	user_pass := ctx.Form.Array("user_pass")
	//利用不可から取得する
	flg_user := ctx.Form.Array("flg_user")
	//チェックから取得する
	flg_exec := ctx.Form.Array("flg_exec")
	//
	flg_update := ctx.Form.String("flg_update")
	//▲▲▲▲▲▲▲▲▲
	//▼▼▼▼▼▼▼▼▼
	//処理
	//============================================
	user := WebApp.GetContextUser(ctx)

	sm := Models.ShopMasterModel{ctx.DB}
	cm := Models.CorpMasterModel{ctx.DB}
	userId := user.UserID
	listShop, err := sm.GetListShopByUser(userId)
	listCorp := []ModelItems.CorpItem{}
	err = cm.GetListCorp(user.CorpCd, &listCorp)
	//=============================================

	var list_user []Models.ListUser
	for i, val := range flg_exec {
		id := strings.TrimSpace(user_id[i])
		if id != "" && val == "1" {
			u := Models.ListUser{}
			u.Um_Shop_Chain_Cd = ""
			u.Um_User_ID = id
			u.Um_User_Name = user_name[i]
			if shop_cd[i] == "-1" {
				shop_cd[i] = ""
				u.Um_Shop_Name = ""
				u.Um_Franchise_Cd = franchise_cd[i]
				u.Um_Franchise_Group_Cd = user.FranchiseGroupCd
			} else if shop_cd[i] == Common.HONBU_SHOP_CD {
				u.Um_Shop_Name = Common.HONBU_SHOP_NAME
				u.Um_Franchise_Cd = ""
				u.Um_Franchise_Group_Cd = ""
			} else {
				u.Um_Franchise_Cd = ""
				u.Um_Franchise_Group_Cd = ""
			}
			u.Um_Shop_Cd = shop_cd[i]
			for _, shop := range listShop {
				if shop.ShopCD == shop_cd[i] {
					u.Um_Shop_Name = shop.ShopName
					u.Um_Server_Name = shop.ServerName
					break
				}
			}

			u.Um_Flg_Auth = flg_auth[i]
			u.Um_Flg_Menu_Group = flg_menu_group[i]
			u.Um_Corp_Cd = corp_cd[i]
			for _, corp := range listCorp {
				if corp.CorpCd == corp_cd[i] {
					u.Um_Corp_Name = corp.CorpName
					break
				}
			}
			u.Um_Dept_Cd = dept_cd[i]
			u.Um_Dept_Name = dept_name[i]
			u.Um_User_Mail = user_mail[i]
			u.Um_User_Phone = user_phone[i]
			u.Um_User_Xerox = user_xerox[i]
			u.Um_User_Pass = user_pass[i]
			u.Um_Flg_Use = flg_user[i]
			list_user = append(list_user, u)
		}
	}
	um := Models.UserMasterModel{ctx.DB}
	err = um.UpdateInsertUser(list_user)
	if err != nil {
		fmt.Println(err)
	}
	//▲▲▲▲▲▲▲▲▲
	//▼▼▼▼▼▼▼▼▼
	//画面
	ctx.ViewData["waiting"] = true
	ctx.ViewData["flg_update"] = flg_update
	ctx.ViewData["listUser"] = list_user
	ctx.ViewData["link_update"] = PATH_MAINTENANCE_USER_SEARCH
	ctx.ViewData["link_insert"] = PATH_MAINTENANCE_USER_NEW

	ctx.View = "maintenance/user/confirm.html"
	//▲▲▲▲▲▲▲▲▲
}
