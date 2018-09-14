package RP063_SingleGoods_Stock_X

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

	item := totaData
	csvWriter.Write([]string{"ＪＡＮ", item.JanCd})
	csvWriter.Write([]string{"品名", item.GoodsName})
	csvWriter.Write([]string{"著者", item.AuthorName})
	csvWriter.Write([]string{"出版社", item.PublisherName})
	csvWriter.Write([]string{"発売日", item.SaleDate})
	csvWriter.Write([]string{"本体価格", fmt.Sprintf("%v", item.Price)})
	csvWriter.Write([]string{"期間入荷累計", fmt.Sprintf("%v", totaData.ReturnTotal)})
	csvWriter.Write([]string{"期間売上累計", fmt.Sprintf("%v", totaData.SaleTotal)})
	csvWriter.Write([]string{"在庫", fmt.Sprintf("%v", totaData.StockTotal)})
	csvWriter.Write([]string{"初売上日", item.FirstSaleDate})

	if csvWriter != nil {
		csvWriter.Write([]string{})
	}

	// 1. detail -------------------------------------------------------

	if groupType == GROUP_TYPE_DATE {
		strRowsArr := []string{}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "年")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcyyyy)
		}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "月")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcmm)
		}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "日")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcdd)
		}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		csvWriter.Write(strRowsArr)
	} else if groupType == GROUP_TYPE_WEEK {
		strRowsArr := []string{}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "年")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcyyyy)
		}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "週")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcweekdate)
		}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		csvWriter.Write(strRowsArr)

	} else {
		strRowsArr := []string{}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "年")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcyyyy)
		}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		csvWriter.Write(strRowsArr)
		strRowsArr = []string{}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "月")
		for _, v := range listRange {
			strRowsArr = append(strRowsArr, v.Mcmm)
		}
		strRowsArr = append(strRowsArr, "")
		strRowsArr = append(strRowsArr, "")
		csvWriter.Write(strRowsArr)
	}

	// ---------------------------------------------------------------
	// 2. Title Row ---------------------------------------------------
	strRowsArrHeadRow := []string{}
	strRowsArrHeadRow = append(strRowsArrHeadRow, "店舗")
	strRowsArrHeadRow = append(strRowsArrHeadRow, "")
	for range listRange {
		strRowsArrHeadRow = append(strRowsArrHeadRow, "合計")
	}
	strRowsArrHeadRow = append(strRowsArrHeadRow, "売上数合計")
	strRowsArrHeadRow = append(strRowsArrHeadRow, "現在在庫数")
	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadRow)
	}

	// ---------------------------------------------------------------
	strRowsArrHeadTotal := []string{}
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "全店合計")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "売上数")
	for _, v := range listRange {
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.SaleDay[v.McKey]))
	}
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.SaleTotalDate))
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.StockCountByShopSearchDate))
	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadTotal)
	}
	strRowsArrHeadTotal = []string{}
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "入荷数")
	for _, v := range listRange {
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, fmt.Sprintf("%v", totaData.ReturnDay[v.McKey]))
	}
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	strRowsArrHeadTotal = append(strRowsArrHeadTotal, "")
	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadTotal)
	}
	//------------------------------------------------------------------

	// 3. Row data
	for _, r := range data {
		strRowsArrRowData := []string{}
		strRowsArrRowData = append(strRowsArrRowData, r.ShopName)
		strRowsArrRowData = append(strRowsArrRowData, "売上数")
		for _, v := range listRange {
			strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.SaleDay[v.McKey]))
		}
		strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.SaleTotalDate))
		strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.StockCountByShop))
		if csvWriter != nil {
			csvWriter.Write(strRowsArrRowData)
		}
		strRowsArrRowData = []string{}
		strRowsArrRowData = append(strRowsArrRowData, "")
		strRowsArrRowData = append(strRowsArrRowData, "入荷数")
		for _, v := range listRange {
			strRowsArrRowData = append(strRowsArrRowData, fmt.Sprintf("%v", r.ReturnDay[v.McKey]))
		}
		strRowsArrRowData = append(strRowsArrRowData, "")
		strRowsArrRowData = append(strRowsArrRowData, "")
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
