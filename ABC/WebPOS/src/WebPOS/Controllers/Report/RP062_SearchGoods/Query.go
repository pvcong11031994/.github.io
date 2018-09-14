package RP062_SearchGoods

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/WebApp"
	"encoding/json"
	"fmt"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"golang.org/x/text/width"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

func Query(ctx *gf.Context) {

	var err error
	ctx.ViewBases = nil
	form := QueryForm{}
	ctx.Form.ReadStruct(&form)
	form.Limit = 1000

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Check アイテム
	if len(form.GoodsType) == 0 {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_SELECTED
		ctx.View = "report/062_search_goods/result_0.html"
		return
	}

	//++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Request HONTO API
	urlRequest := ctx.Config.StrOrEmpty(WebApp.CONFIG_SEARCH_GOODS_URL) + "?"
	values := url.Values{}
	access_key := ctx.Config.StrOrEmpty(WebApp.CONFIG_ACCESS_KEY_HONTO_API)
	if access_key == "" {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_SYSTEM
		ctx.View = "report/062_search_goods/result_0.html"
		return
	}
	values.Set("access_key", access_key)
	//Check len keyword > 2
	flagCheckLen := CheckLen([]byte(form.KeyWord))
	if flagCheckLen {
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_LEN_KEY_WORD
		ctx.View = "report/062_search_goods/result_0.html"
		return
	}
	//Check keyword
	if len(form.KeyWord) > 0 {
		//Split キーワード
		form.KeyWordArrays = strings.Split(form.KeyWord, " ")
		for _, value := range form.KeyWordArrays {
			values.Add("keyword", value)
		}
	}

	//Check Item checkbox
	if len(form.GoodsType) == 1 {
		if form.GoodsType[0] == "1" {
			values.Set("item_name", "本")
		} else if form.GoodsType[0] == "2" {
			values.Set("item_name", "雑誌")
		}
	}

	//Check sort by
	values.Set("sort", form.Sort)
	//Add count
	values.Set("count", TOTAL_COUNT_PAGE)
	//Check page
	values.Set("page", form.Page)
	urlRequest += values.Encode()

	/*proxyUrl, err := url.Parse("http://proxy.fujinet.vn:8080")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
	}
	client := &http.Client{Transport: tr}*/

	//Get Request
	resp, err := http.Get(urlRequest)
	if err != nil {
		Common.LogErr(err)
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_SYSTEM
		ctx.View = "report/062_search_goods/result_0.html"
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	listRsp := Response{}

	jsonErr := json.Unmarshal(body, &listRsp)
	if jsonErr != nil {
		Common.LogErr(err)
		ctx.ViewData["err_msg"] = RPComon.REPORT_ERROR_SYSTEM
		ctx.View = "report/062_search_goods/result_0.html"
		return
	}
	// set report name to import info log search charging
	ctx.SetSessionFlash(RPComon.REPORT_NAME_KEY, _REPORT_NAME)
	// ========================================================================================
	// Output log search condition
	tag := "report=" + _REPORT_ID
	if len(form.KeyWord) > 0 {
		tag += ",キーワード=" + `"` + form.KeyWord + `"`
	}
	if len(form.GoodsType) == 1 {
		value := ""
		if form.GoodsType[0] == "1" {
			value = "本"
		} else if form.GoodsType[0] == "2" {
			value = "雑誌"
		}
		tag += ",アイテム=" + `"` + value + `"`
	}
	tag += ",並び順=" + `"` + form.Sort + `"`
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
	tag = tag + `,app_id="mBAWEB-v13a"`
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END

	queryLog := bq.QueryLog{
		Context:   ctx,
		Tag:       tag,
		StartAt:   time.Now(),
		QuerySize: 0,
		ExecTime:  0,
		State:     bq.QUERY_LOG_BEGIN,
	}
	RPComon.QueryLogHandle(&queryLog)
	queryLog = bq.QueryLog{
		Context:   ctx,
		Tag:       tag,
		StartAt:   time.Now(),
		QuerySize: 0,
		ExecTime:  0,
		State:     bq.QUERY_LOG_END,
	}
	RPComon.QueryLogHandle(&queryLog)

	//-------------------------------------------------------------
	//Convert date
	for i, valueDate := range listRsp.ProductList {
		listRsp.ProductList[i].Release.ReleaseDate = FormatDateTime(valueDate.Release.ReleaseDate)
	}
	ctx.ViewData["data"] = listRsp
	ctx.ViewData["form"] = form
	ctx.ViewData["listJan"] = form.JanArrays
	ctx.TemplateFunc["arr"] = Common.MakeArray
	ctx.TemplateFunc["convert_int"] = Common.ConvertStringToInt
	ctx.TemplateFunc["sum_format"] = Common.FormatNumber
	ctx.TemplateFunc["date_format"] = Common.FormatDateTime

	if listRsp.TotalItemCount == "0" {
		msgError := RPComon.REPORT_ERROR_SYSTEM
		if listRsp.Status == "0000" {
			msgError = RPComon.REPORT_SEARCH_RESULT_EMPTY
		}
		ctx.ViewData["err_msg"] = msgError
		ctx.View = "report/062_search_goods/result_0.html"
		return
	} else {
		ctx.View = "report/062_search_goods/result_4.html"
	}
}

func FormatDateTime(strDateTime string) string {

	var resultTempReleaseDate = ""
	if len(strDateTime) > 0 {
		arrTempReleaseDate := strings.Split(strDateTime, ":")
		arrTempReleaseDate2Byte := strings.Split(strDateTime, "：")
		if len(arrTempReleaseDate) > len(arrTempReleaseDate2Byte) {
			resultTempReleaseDate = arrTempReleaseDate[1]
		} else if len(arrTempReleaseDate) < len(arrTempReleaseDate2Byte) {
			resultTempReleaseDate = arrTempReleaseDate2Byte[1]
		} else {
			resultTempReleaseDate = arrTempReleaseDate[0]
		}
	}

	//Convert 2byte to 1byte
	result := width.Narrow.String(resultTempReleaseDate)
	fmt.Sprintf("%s", result)

	return result
}

func CheckLen(strConvert []byte) bool {

	flag := 0
	for len(strConvert) > 0 {
		flag++
		_, size := utf8.DecodeLastRune(strConvert)
		strConvert = strConvert[:len(strConvert)-size]
	}
	if flag <= 1 {
		return true
	} else {
		return false
	}
}
