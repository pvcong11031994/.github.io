package RPComon

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"WebPOS/WebApp"
	"github.com/goframework/gcp/bq"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"log"
	"math"
	"strings"
)

func init() {
	bq.SetDefaultQueryLogHandler(QueryLogHandle)
}

func QueryLogHandle(queryLog *bq.QueryLog) {
	if queryLog != nil {
		if queryLog.Context != nil {
			ctx, ctxOK := queryLog.Context.(*gf.Context)
			if ctxOK {
				user := WebApp.GetContextUser(ctx)
				switch queryLog.State {
				case bq.QUERY_LOG_BEGIN:
					//log.Printf("QUERY LOG BEGIN: AppVersion=%v, user=%v,report=%v", WebApp.APP_VERSION, user.UserName, queryLog.Tag)
					log.Printf("QUERY LOG BEGIN: AppVersion=%v, user=%v,%v", WebApp.APP_VERSION, user.UserName, queryLog.Tag)
				case bq.QUERY_LOG_END:
					//log.Printf("QUERY LOG END: AppVersion=%v, user=%v,report=%v,Size=%v * %v,Time=%v\r\n", WebApp.APP_VERSION, user.UserName, queryLog.Tag, queryLog.QuerySize, queryLog.BillingTier, queryLog.ExecTime)
					size := Common.ByteSize(uint64(queryLog.QuerySize * queryLog.BillingTier))
					prmArr := strings.Split(queryLog.Tag, ",")
					strReport := ""
					if len(prmArr) > 0 {
						strReport = prmArr[0]
					}
					log.Printf("QUERY LOG END: AppVersion=%v, user=%v,%v,Size=%v,Time=%v\r\n", WebApp.APP_VERSION, user.UserName, strReport, size, queryLog.ExecTime)
					err := insertLogChargingHandle(queryLog)
					Common.LogErr(err)
				}
			}
		}
	}
}

// write info search into table bq_log_search_charging
func insertLogChargingHandle(queryLog *bq.QueryLog) error {

	ctx, ctxOK := queryLog.Context.(*gf.Context)
	if !ctxOK {
		return nil
	}

	// get info user
	user := WebApp.GetContextUser(ctx)

	// get info shop
	shm := Models.ShopMasterModel{ctx.DB}
	shopInfo, err := shm.GetInfoShopByCD(user.ShopCd)
	if err != nil {
		return exterror.WrapExtError(err)
	}

	// get GCP_charging
	var floatGcpCharging float64 = 0
	floatGcpCharging = float64(REPORT_CHARGING_GOOGLE_API_QUERY*REPORT_EXCHANGE_RATE*queryLog.QuerySize*queryLog.BillingTier) / math.Pow(1024, 4)
	intGcpCharging := int(math.Ceil(floatGcpCharging * 1000))

	// get bqcs_bairitsu from table bq_charging_setting
	bqcs := Models.BQChargingSettingModel{ctx.DB}
	infoCharging, err := bqcs.GetInfoChargingByCd(user.ShopCd, user.ServerName, user.UserID)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	// get VJ_charging
	var vjCharging int = 0
	if infoCharging.ShopCd != "" {
		vjCharging = intGcpCharging * infoCharging.Bairitsu
	}
	// return value vj_charging to display on result screen
	ctx.SetSessionFlash(REPORT_VJ_CHARGING_KEY, vjCharging)

	// get report name
	menuReportName := ""
	if reportName := ctx.GetSessionFlash(REPORT_NAME_KEY); reportName != nil {
		menuReportName = reportName.(string)
	}

	// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5719 - ADD - START
	strHandle := ""
	strFormat := ""
	strTab := ""
	strAppID := ""

	tag := queryLog.Tag
	tagArr := strings.Split(tag, ",")
	if len(tagArr) > 0 {
		for _, v := range tagArr {
			v = strings.Replace(strings.Replace(v, " ", "", -1), `"`, "", -1)
			if strings.Contains(v, "handle") {
				strHandle = strings.Replace(v, "handle=", "", -1)
			}
			if strings.Contains(v, "フォーマット") {
				strFormat = strings.Replace(v, "フォーマット=", "", -1)
			}
			if strings.Contains(v, "tab") {
				strTab = strings.Replace(v, "tab=", "", -1)
			}
			if strings.Contains(v, "app_id") {
				strAppID = strings.Replace(v, "app_id=", "", -1)
			}
		}
	}
	// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5719 - ADD - END

	itemInsert := ModelItems.BQLogSearchChargingItem{
		ServerName:    user.ServerName,
		ShopCd:        user.ShopCd,
		ShopName:      shopInfo.ShopName,
		FranchiseCd:   user.FranchiseCd,
		UserID:        user.UserID,
		UserName:      user.UserName,
		UserMenu:      menuReportName,
		UseTBLSize:    queryLog.QuerySize * queryLog.BillingTier,
		GCPCharging:   intGcpCharging,
		VJCharging:    vjCharging,
		ExecTime:      queryLog.ExecTime,
		AppVersion:    WebApp.APP_VERSION,
		SearchingDate: queryLog.StartAt.Format(Common.DATE_FORMAT_MYSQL_YMDHMS),
		Handle:        strHandle,
		Format:        strFormat,
		Tab:           strTab,
		AppID:         strAppID,
	}

	bqlsc := Models.BQLogSearchChargingModel{ctx.DB}
	err = bqlsc.InsertLogChargingByUser(itemInsert)
	if err != nil {
		return exterror.WrapExtError(err)
	}

	return nil
}
