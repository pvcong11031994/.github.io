package RP060_FavoriteManagement

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"encoding/json"
	"fmt"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strings"
)

type UserJan struct {
	JanCode        string
	PriorityNumber string
	ProductName    string
	AuthorName     string
	MakerName      string
	UnitPrice      string
	ReleaseDate    string
	Memo           string
}

// 初速比較検索
func Load(ctx *gf.Context) {

	ctx.ViewBases = nil
	// Get info user
	user := WebApp.GetContextUser(ctx)

	// Get List shop by user role
	ujc := Models.NewUserJanCodeModel(ctx.DB)
	listUserJan, err := ujc.GetListUserJanByUser(user.UserID)
	Common.LogErr(err)

	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.ViewData["list_user_jan"] = listUserJan
	ctx.View = "report/060_favorite_management/result_4.html"
}

//Check validate and show data to screen
func Delete(ctx *gf.Context) {

	var err error
	janList := ctx.Form.String("list_jan")
	janList = janList[:len(janList)-1]
	janListArray := strings.Split(janList, ";")
	user := WebApp.GetContextUser(ctx)
	ujc := Models.NewUserJanCodeModel(ctx.DB)
	err = ujc.DeleteListUserJan(user.UserID, janListArray)
	if err == nil {
		ctx.JsonResponse = map[string]interface{}{
			"Success": true,
			"Msg":     RPComon.REPORT_DELETE_FAVORITE_SUCCESS,
		}
	} else {
		Common.LogErr(err)
		ctx.JsonResponse = map[string]interface{}{"Success": false, "Msg": RPComon.REPORT_DELETE_FAVORITE_FAIL}
	}
	for _, jan := range janListArray {
		Common.LogOutput(fmt.Sprintf(Common.FAVORITE_MANAGEMENT_LOG, user.UserID, jan, Common.FAVORITE_MANAGEMENT_DELETE_PROCESS))
	}
}

func Update(ctx *gf.Context) {

	var err error
	data := ctx.Form.StringNoEscape("data")
	bytes := []byte(data)
	var userJanData []UserJan
	json.Unmarshal(bytes, &userJanData)

	user := WebApp.GetContextUser(ctx)
	ujc := Models.NewUserJanCodeModel(ctx.DB)
	err = ujc.BeginTrans()
	if err != nil {
		exterror.WrapExtError(err)
	}
	for _, userJan := range userJanData {
		args := []interface{}{}
		args = append(args, userJan.PriorityNumber, userJan.Memo, user.UserID, userJan.JanCode)
		err = ujc.UpdateInfoListUserJan(args)
		if err != nil {
			break
		}
		Common.LogOutput(fmt.Sprintf(Common.FAVORITE_MANAGEMENT_LOG, user.UserID, userJan.JanCode, Common.FAVORITE_MANAGEMENT_UPDATE_PROCESS))
	}
	ujc.FinishTrans(&err)
	ujc.DeleteAutoUserJan(user.UserID)
	if err == nil {
		ctx.JsonResponse = map[string]interface{}{
			"Success": true,
			"Msg":     RPComon.REPORT_UPDATE_FAVORITE_SUCCESS,
		}
	} else {
		Common.LogErr(err)
		ctx.JsonResponse = map[string]interface{}{"Success": false, "Msg": RPComon.REPORT_UPDATE_FAVORITE_FAIL}
	}
}

func Insert(ctx *gf.Context) {

	var err error
	var janError string
	data := ctx.Form.StringNoEscape("data")
	bytes := []byte(data)
	var userJanData []UserJan
	json.Unmarshal(bytes, &userJanData)

	user := WebApp.GetContextUser(ctx)
	ujc := Models.NewUserJanCodeModel(ctx.DB)
	err = ujc.BeginTrans()
	if err != nil {
		exterror.WrapExtError(err)
	}
	for _, userJan := range userJanData {
		if len(userJan.JanCode) == 13 {
			args := []interface{}{}
			// ASO-5873 mBAWEB-v14a お気に入り管理：Maria⇒CLOUDSQL引越 - ADD START
			if strings.Compare(userJan.UnitPrice, "") == 0 {
				userJan.UnitPrice = "0"
			}
			// ASO-5873 mBAWEB-v14a お気に入り管理：Maria⇒CLOUDSQL引越 - ADD END
			args = append(args, user.UserID, userJan.JanCode, userJan.ProductName, userJan.MakerName,
				userJan.AuthorName, userJan.ReleaseDate, userJan.UnitPrice, userJan.Memo, userJan.PriorityNumber)
			err = ujc.InsertUpdateUserJan(args)
			if err != nil {
				janError = userJan.JanCode
				break
			}
			Common.LogOutput(fmt.Sprintf(Common.FAVORITE_MANAGEMENT_LOG, user.UserID, userJan.JanCode, Common.FAVORITE_MANAGEMENT_UPDATE_PROCESS))
		}
	}
	ujc.FinishTrans(&err)
	ujc.DeleteAutoUserJan(user.UserID)
	if err == nil {
		ctx.JsonResponse = map[string]interface{}{
			"Success": true,
			"Msg":     RPComon.REPORT_INSERT_FAVORITE_SUCCESS,
		}
	} else {
		Common.LogErr(err)
		ctx.JsonResponse = map[string]interface{}{
			"Success":  false,
			"Msg":      RPComon.REPORT_INSERT_FAVORITE_FAIL,
			"JanError": janError,
		}
	}
}

func UpdateJanFavorite(ctx *gf.Context, userJan UserJan) (err error) {
	user := WebApp.GetContextUser(ctx)
	ujc := Models.NewUserJanCodeModel(ctx.DB)
	err = ujc.BeginTrans()
	if err != nil {
		exterror.WrapExtError(err)
	}
	if len(userJan.JanCode) == 13 {
		args := []interface{}{}
		args = append(args, user.UserID, userJan.JanCode, userJan.ProductName, userJan.MakerName,
			userJan.AuthorName, userJan.ReleaseDate, userJan.UnitPrice, userJan.Memo, userJan.PriorityNumber)
		err = ujc.InsertUpdateUserJan(args)
		ujc.FinishTrans(&err)
		ujc.DeleteAutoUserJan(user.UserID)
		Common.LogOutput(fmt.Sprintf(Common.FAVORITE_MANAGEMENT_LOG, user.UserID, userJan.JanCode, Common.FAVORITE_MANAGEMENT_UPDATE_PROCESS))
	}
	return err
}

//Insert data when link into search_goods
func InsertOrUpdate(ctx *gf.Context, form QueryForm) {

	var err error
	var janError string

	//Cut Array Data
	listJanSelected := strings.Split(form.JanCodeList, DEFAULT_KEY_LINK_UPDATE_INSERT)
	listProductNameSelected := strings.Split(form.ProductNameList, DEFAULT_KEY_LINK_UPDATE_INSERT)
	listAuthorNameSelected := strings.Split(form.AuthorNameList, DEFAULT_KEY_LINK_UPDATE_INSERT)
	listMakerNameSelected := strings.Split(form.MakerNameList, DEFAULT_KEY_LINK_UPDATE_INSERT)
	listUnitPriceListSelected := strings.Split(form.UnitPriceList, DEFAULT_KEY_LINK_UPDATE_INSERT)
	listReleaseDateSelected := strings.Split(form.ReleaseDateList, DEFAULT_KEY_LINK_UPDATE_INSERT)

	user := WebApp.GetContextUser(ctx)
	model := Models.NewUserJanCodeModel(ctx.DB)
	err = model.BeginTrans()
	for i := 0; i < form.LengthListSelected; i++ {
		if len(listJanSelected[i]) == 13 {
			args := []interface{}{}
			listUnitPriceListSelected[i] = FormatPrice(listUnitPriceListSelected[i])
			args = append(args, user.UserID, listJanSelected[i], listProductNameSelected[i], listMakerNameSelected[i],
				listAuthorNameSelected[i], listReleaseDateSelected[i], listUnitPriceListSelected[i], "", DEFAULT_PRIORITY_NUMBER)
			err = model.InsertUpdateUserJanSearchGoods(args)
			if err != nil {
				janError = listJanSelected[i]
				break
			}
			Common.LogOutput(fmt.Sprintf(Common.FAVORITE_MANAGEMENT_LOG, user.UserID, listJanSelected[i], Common.FAVORITE_MANAGEMENT_UPDATE_PROCESS))
		}
	}
	model.FinishTrans(&err)
	model.DeleteAutoUserJan(user.UserID)
	if err != nil {
		Common.LogErr(err)
		ctx.JsonResponse = map[string]interface{}{
			"Success":  false,
			"Msg":      RPComon.REPORT_INSERT_FAVORITE_FAIL,
			"JanError": janError,
		}
	}
}

func FormatPrice(strPrice string) string {
	strResult := strings.Replace(strPrice, ",", "", -1)
	return strResult
}
