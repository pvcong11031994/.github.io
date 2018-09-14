package RP057_SingleGoods_Cumulative

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"WebPOS/Models/ModelItems"
	"encoding/csv"
	"fmt"
	"github.com/goframework/encode"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"runtime"
)

/*
	4. Write file ROW | COL | SUM
*/

func WriteFile4SingleNew(data []SingleItem, totaData SingleItem, searchHandleType, groupType string, listRange []ModelItems.MasterCalendarItem) (error, string, string) {

	defer runtime.GC()
	var csvWriter *csv.Writer = nil
	var excelWriter *Common.SimpleExcelFile = nil
	var csvFile *os.File = nil

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	fileName := _REPORT_ID + "_" + Common.CurrentDateTime()
	filePath := filepath.Join(fileDir, fileName)
	if searchHandleType == RPComon.REPORT_SEARCH_TYPE_HANDLE_CSV {
		filePath = filePath + ".csv"
		fileName = fileName + ".csv"
		csvFile, err := os.Create(filePath)
		if err != nil {
			return exterror.WrapExtError(err), "", ""
		}
		csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(csvFile))
		defer csvFile.Close()
		csvWriter.UseCRLF = true
	}

	// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5559#comment-3260016
	//item := data[0]
	item := data[len(data)-1]
	csvWriter.Write([]string{"ＪＡＮ", item.JanCd})
	csvWriter.Write([]string{"品名", item.GoodsName})
	csvWriter.Write([]string{"著者", item.AuthorName})
	csvWriter.Write([]string{"出版社", item.PublisherName})
	csvWriter.Write([]string{"発売日", item.SaleDate})
	csvWriter.Write([]string{"本体価格", fmt.Sprintf("%v", item.Price)})
	// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5559#comment-3260016
	//csvWriter.Write([]string{"入荷累計", fmt.Sprintf("%v", totaData.ReturnTotal)})
	//csvWriter.Write([]string{"売上累計", fmt.Sprintf("%v", totaData.SaleTotal)})
	//csvWriter.Write([]string{"在庫", fmt.Sprintf("%v", totaData.StockCurCount)})
	csvWriter.Write([]string{"入荷累計", fmt.Sprintf("%v", item.ReturnTotal)})
	csvWriter.Write([]string{"売上累計", fmt.Sprintf("%v", item.SaleTotal)})
	csvWriter.Write([]string{"在庫", fmt.Sprintf("%v", item.StockCurCount)})
	csvWriter.Write([]string{"初売上日", item.FirstSaleDate})
	// ASO-5651 [BA]mBAWEB-v02a 単品推移：商品情報に出版社在庫の表示を追加
	// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5597#comment-3267291
	csvWriter.Write([]string{"出版社在庫", item.StockInf})

	if csvWriter != nil {
		csvWriter.Write([]string{})
	}

	// 1. detail -------------------------------------------------------

	if groupType == GROUP_TYPE_DATE {
		strRowsArr := []string{}
		strRowsArr = append(strRowsArr, "")
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "年")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcyyyy)
		}
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "月")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcmm)
		}
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "日")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcdd)
		}
		csvWriter.Write(strRowsArr)
	} else if groupType == GROUP_TYPE_WEEK {
		strRowsArr := []string{}
		strRowsArr = append(strRowsArr, "")
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "年")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcyyyy)
		}
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "週")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcweekdate)
		}
		csvWriter.Write(strRowsArr)

	} else {
		strRowsArr := []string{}
		strRowsArr = append(strRowsArr, "")
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "年")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcyyyy)
		}
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "月")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcmm)
		}
		csvWriter.Write(strRowsArr)
	}

	// ---------------------------------------------------------------
	// 2. Title Row ---------------------------------------------------
	strRowsArrHeadRow := []string{}
	// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
	strRowsArrHeadRow = append(strRowsArrHeadRow, "共有書店コード")
	strRowsArrHeadRow = append(strRowsArrHeadRow, "店舗コード")
	strRowsArrHeadRow = append(strRowsArrHeadRow, "店舗")
	// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
	strRowsArrHeadRow = append(strRowsArrHeadRow, "在庫数")
	strRowsArrHeadRow = append(strRowsArrHeadRow, "期間合計")
	strRowsArrHeadRow = append(strRowsArrHeadRow, "")
	for range listRange {
		strRowsArrHeadRow = append(strRowsArrHeadRow, "合計")
	}
	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadRow)
	}

	// ---------------------------------------------------------------
	strRowsArrHeadTotal := []string{}
	// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "全店合計")
	// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.StockCurCount))
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.ReturnTotalDate))
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "入荷数")
	for _, v := range listRange {
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.ReturnDay[v.McKey]))
	}
	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadTotal)
	}
	strRowsArrHeadTotal = []string{}
	// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "全店合計")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.SaleTotalDate))
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "売上数")
	for _, v := range listRange {
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.SaleDay[v.McKey]))
	}
	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadTotal)
	}
	//------------------------------------------------------------------

	// 3. Row data
	for i, r := range data {
		// https://dev-backlog.rsp.honto.jp/backlog/view/ASO-5559#comment-3260016
		if i >= len(data) - 1 {
			break
		}
		strRowsArrRowData := []string{}
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5597#comment-3267291
		//strRowsArrRowData = append(strRowsArrRowData, r.ShopSeqNumber)
		strRowsArrRowData = append(strRowsArrRowData, r.SharedBookStoreCode)
		strRowsArrRowData = append(strRowsArrRowData, r.ShopCd)
		strRowsArrRowData = append(strRowsArrRowData, r.ShopName)
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.StockCurCount))
		strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.ReturnTotalDate))
		strRowsArrRowData = append(strRowsArrRowData, "入荷数")
		for _, v := range listRange {
			strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.ReturnDay[v.McKey]))
		}
		if csvWriter != nil {
			csvWriter.Write(strRowsArrRowData)
		}
		strRowsArrRowData = []string{}
		// ASO-5597 [BA]mBAWEB-v02a 単品推移：商品情報のCLOUD化と店舗マスタの利用
		// https://backlog-dev.rsp.honto.jp/backlog/view/ASO-5597#comment-3267291
		//strRowsArrRowData = append(strRowsArrRowData, r.ShopSeqNumber)
		strRowsArrRowData = append(strRowsArrRowData, r.SharedBookStoreCode)
		strRowsArrRowData = append(strRowsArrRowData, r.ShopCd)
		strRowsArrRowData = append(strRowsArrRowData, r.ShopName)
		strRowsArrRowData = append(strRowsArrRowData, "")
		strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.SaleTotalDate))
		strRowsArrRowData = append(strRowsArrRowData, "売上数")
		for _, v := range listRange {
			strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.SaleDay[v.McKey]))
		}
		if csvWriter != nil {
			csvWriter.Write(strRowsArrRowData)
		}
	}

	// ---------------------------------------------------------------
	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}
	if excelWriter != nil {
		excelWriter.Close()
	}

	return nil, filePath, fileName
}
