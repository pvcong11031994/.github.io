package UsersMaintenance

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
)

func Search(ctx *gf.Context) {
	sm := Models.ShopMasterModel{ctx.DB}
	user := WebApp.GetContextUser(ctx)

	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)
	ctx.ViewData["shops"] = listShop

	ctx.ViewData["honbu_shop_cd"] = Common.HONBU_SHOP_CD
	ctx.ViewData["honbu_shop_name"] = Common.HONBU_SHOP_NAME
	ctx.ViewData["flg_auth"] = user.FlgAuth

	ctx.ViewData["link_result_list"] = PATH_MAINTENANCE_USER_SEARCH_LIST
	ctx.View = "maintenance/user/search.html"
}

func List(ctx *gf.Context) {

	// 店舗から取得する
	lShop := ctx.Form.Array("shop_cd")
	// ユーザ名から取得する
	user_name := ctx.Form.String("user_name")
	// 権限から取得する
	flg_auth := ctx.Form.String("flg_auth")

	user := WebApp.GetContextUser(ctx)

	um := Models.UserMasterModel{ctx.DB}
	sm := Models.ShopMasterModel{ctx.DB}
	cm := Models.CorpMasterModel{ctx.DB}

	listUser, err := um.GetListUser(lShop, user_name, flg_auth)
	Common.LogErr(err)

	listShop, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	listMenu := um.ListMenuGroup(user.ShopChainCd, user.ShopCd)
	sfgm := Models.ShopFranchiseGroupMasterModel{ctx.DB}

	listFranchiseCd := []ModelItems.FranchiseCdItem{}
	err = sfgm.GetListFranchiseCD(user.FranchiseGroupCd, &listFranchiseCd)

	listCorp := []ModelItems.CorpItem{}
	err = cm.GetListCorp(user.CorpCd, &listCorp)
	Common.LogErr(err)

	ctx.ViewData["corps"] = listCorp
	ctx.ViewData["shops"] = listShop
	ctx.ViewData["menugroups"] = listMenu
	ctx.ViewData["listUser"] = listUser
	ctx.ViewData["listFranchiseCd"] = listFranchiseCd

	ctx.ViewData["honbu_shop_cd"] = Common.HONBU_SHOP_CD
	ctx.ViewData["honbu_shop_name"] = Common.HONBU_SHOP_NAME
	ctx.ViewData["shop_chain"] = user.ShopChainCd

	ctx.ViewData["link_search"] = PATH_MAINTENANCE_USER_SEARCH
	ctx.ViewData["link_result_confirm"] = PATH_MAINTENANCE_USER_CONFIRM
	ctx.View = "maintenance/user/list.html"
}
