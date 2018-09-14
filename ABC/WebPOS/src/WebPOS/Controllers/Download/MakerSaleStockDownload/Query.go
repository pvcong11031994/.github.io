package MakerSaleStockDownload

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/DB"
	"WebPOS/WebApp"
	"encoding/csv"
	"github.com/goframework/encode"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
)

const (
	_MSS_DATA_DOWNLOAD_KEY_ = "_MSS_DATA_DOWNLOAD_KEY_"
	// FUJI-4677 EDIT
	//_ERROR_NO_DATA          = "データがありません。<br/>条件を変更して再度検索してください。"
	_ERROR_NO_DATA          = "該当するデータがありません。<br/>条件を変更して再度検索してください。"
	_ERROR_MAKER            = "出版社を入力してください。"
	_ERROR_DATE             = "日付を入力してください。"
)

func Query(ctx *gf.Context) {

	var err error
	ctx.ViewBases = nil
	user := WebApp.GetContextUser(ctx)
	form := Form{}
	ctx.Form.ReadStruct(&form)

	// 同時処理管理処理
	mPoolChan <- true
	defer func() { <-mPoolChan }()

	// 出版社とJAN入力チェック++++++++++++++++++++++++++++++++++++++++++++++
	if len(form.MakerCd) == 0 && len(form.JAN) == 0 {
		onError(ctx, _ERROR_MAKER)
		return
	}
	//==========================================================

	sm := Models.ShopMasterModel{ctx.DB}
	selectableShops, err := sm.GetListShopByUser(user.UserID)
	Common.LogErr(err)

	// 店舗チェック+++++++++++++++++++++++++++++++++++++++++++++
	selectedShopCd := []string{}
	for _, shopCd := range form.ShopCd {
		for _, item := range selectableShops {
			if shopCd == item.ServerName+"|"+item.ShopCD {
				selectedShopCd = append(selectedShopCd, shopCd)
				break
			}
		}
	}

	if len(selectedShopCd) == 0 {
		onError(ctx, RPComon.REPORT_ERROR_NO_SHOP)
		return
	} else {
		form.ShopCd = selectedShopCd
	}
	//==========================================================

	// 日付チェック+++++++++++++++++++++++++++++++++++++++++++++
	selectedDateFrom := form.DateFrom
	if !Common.IsValidateDate(selectedDateFrom) {
		onError(ctx, _ERROR_DATE)
		return
	}
	selectedDateTo := form.DateTo
	if !Common.IsValidateDate(selectedDateTo) {
		onError(ctx, _ERROR_DATE)
		return
	}
	//==========================================================
	// JAN整形+++++++++++++++++++++++++++++++++++++++++++++++++
	form.JAN = Common.GenerateJAN(form.JAN)
	//==========================================================

	mssm := Models.MakerSaleStockModel{ctx.DB}
	searchResultChan := mssm.Search(form.MakerCd, form.DateFrom, form.DateTo, form.JAN, form.DataMode, form.GoodsType, form.ShopCd, ctx)

	// CSVファイルをエクスポートする
	err, handle, filePath, strErr := writeOutputFile(searchResultChan, form.DataMode)
	if err != nil || filePath == "" {
		Common.LogErr(err)
		if handle {
			strErr = RPComon.REPORT_ERROR_SYSTEM
		}
		onError(ctx, strErr)
		return
	}

	downloadToken := Common.RandString(8)
	ctx.Session.Values[_MSS_DATA_DOWNLOAD_KEY_+downloadToken] = filePath
	ctx.JsonResponse = map[string]interface{}{
		"Success":       true,
		"Msg":           strErr,
		"DownloadToken": downloadToken,
	}
}

func onError(ctx *gf.Context, err string) {

	ctx.JsonResponse = map[string]interface{}{
		"Success": false,
		"Msg":     err,
	}
}

func writeOutputFile(mssDataChan chan Models.MakerSSDataWithErr, dataMode string) (error, bool, string, string) {

	countRecord := 0
	headerWritten := false

	var csvWriter *csv.Writer = nil
	var csvFile *os.File = nil
	var err error

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	filePath := filepath.Join(fileDir, "maker_sale_stock_"+Common.CurrentDateTime()+".csv")
	csvFile, err = os.Create(filePath)
	if err != nil {
		return exterror.WrapExtError(err), true, "", ""
	}
	csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(csvFile))
	csvWriter.UseCRLF = true

	for {
		dataChan, ok := <-mssDataChan
		if !ok {
			break
		}
		if dataChan.Err != nil {
			if csvFile != nil {
				csvFile.Close()
			}
			os.Remove(fileDir)
			return exterror.WrapExtError(dataChan.Err), true, "", ""
		}

		// データモードにより欄を判断する
		if !headerWritten {
			if csvWriter != nil {
				header := Common.ListTagByKey(dataChan.Data, "header")
				switch dataMode {
				case "0":
					header = append(header, "売上数")
				case "1":
					header = append(header, "在庫数")
				}
				csvWriter.Write(header)
			}
			headerWritten = true
		}

		if csvWriter != nil {
			line := Common.ListStringValueByTag(dataChan.Data, "header")
			switch dataMode {
			case "0":
				line = append(line, dataChan.Data.GoodsCount)
			case "1":
				line = append(line, dataChan.Data.StockCount)
			}
			csvWriter.Write(line)
		}
		countRecord += 1
	}

	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}
	if countRecord == 0 {
		return nil, false, "", _ERROR_NO_DATA
	}
	return nil, false, filePath, ""
}
