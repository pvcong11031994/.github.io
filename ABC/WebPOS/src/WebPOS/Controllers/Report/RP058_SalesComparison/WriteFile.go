package RP058_SalesComparison

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"encoding/csv"
	"github.com/goframework/encode"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func WriteFile4(dataSingleItem []SingleItem, searchHandleType, groupType string, ctx *gf.Context, form QueryForm) (error, string, string) {

	defer runtime.GC()
	var csvWriter *csv.Writer = nil
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
	// 1. Head -------------------------------------------------------
	headerRows := []string{
		"",
		"ＪＡＮ",
		"品名",
		"著者",
		"出版社名",
		"発売日",
		"本体価格",
		"入荷累計",
		"売上累計",
		"在庫数",
		"初売上日",
		"期間売上合計",
	}
	//Write header
	strRowsArr := []string{}
	for i := 0; i < len(headerRows); i++ {
		strRowsArr = append(strRowsArr, headerRows[i])
	}
	if csvWriter != nil {
		csvWriter.Write(strRowsArr)
	}
	// 2. Data -------------------------------------------------------
	for i, data := range dataSingleItem {
		if i == RPComon.REPORT_CSV_ROW_LIMIT {
			break
		}
		strData := []string{}
		switch i {
		case 0:
			strData = append(strData, "A")
		case 1:
			strData = append(strData, "B")
		case 2:
			strData = append(strData, "C")
		case 3:
			strData = append(strData, "D")
		case 4:
			strData = append(strData, "E")
		case 5:
			strData = append(strData, "F")
		case 6:
			strData = append(strData, "G")
		case 7:
			strData = append(strData, "H")
		case 8:
			strData = append(strData, "I")
		case 9:
			strData = append(strData, "J")
		case 10:
			strData = append(strData, "K")
		case 11:
			strData = append(strData, "L")
		case 12:
			strData = append(strData, "M")
		case 13:
			strData = append(strData, "N")
		case 14:
			strData = append(strData, "O")
		}

		strData = append(strData, data.JanCd)
		strData = append(strData, data.GoodsName)
		strData = append(strData, data.AuthorName)
		strData = append(strData, data.PublisherName)
		strData = append(strData, data.SaleDate)
		strData = append(strData, strconv.Itoa(int(data.Price)))
		strData = append(strData, strconv.Itoa(int(data.ReturnTotal)))
		strData = append(strData, strconv.Itoa(int(data.SaleTotal)))
		strData = append(strData, strconv.Itoa(int(data.StockCurCount)))
		strData = append(strData, data.FirstSaleDate)
		strData = append(strData, strconv.Itoa(int(data.SaleTotalDate)))
		if csvWriter != nil {
			csvWriter.Write(strData)
		}
	}

	//Write info detail Jan
	for _, data := range dataSingleItem {
		//Call QueryDetail
		dataSingleItemJan := []SingleItem{}
		totalSingleItemJan := SingleItem{}
		sql, listRange := buildDetailSql(form, ctx, data.JanCd)
		dataSingleItemNew, totalSingleItemNew, err := queryDataDetail(ctx, sql, listRange, form, data.JanCd)
		Common.LogErr(err)
		// システムエラー
		if err != nil {
			ctx.ViewData[RPComon.REPORT_ERROR_SYSTEM_VIEW] = RPComon.REPORT_ERROR_SYSTEM
			ctx.View = RPComon.REPORT_ERROR_PATH_HTML
			return nil, filePath, fileName
		}
		dataSingleItemJan = dataSingleItemNew
		totalSingleItemJan = totalSingleItemNew

		//Write File JAN:
		csvWriter.Write([]string{""})
		csvWriter.Write([]string{"JAN", data.JanCd})
		if groupType == GROUP_TYPE_DATE {
			//Year
			strRowsArrJan := []string{}
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "年")
			for _, v := range listRange {
				strRowsArrJan = append(strRowsArrJan, v.Mcyyyy)
			}
			csvWriter.Write(strRowsArrJan)

			//Month
			strRowsArrJan = []string{}
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "月")
			for _, v := range listRange {
				strRowsArrJan = append(strRowsArrJan, v.Mcmm)
			}
			csvWriter.Write(strRowsArrJan)

			//Day
			strRowsArrJan = []string{}
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "日")
			for _, v := range listRange {
				strRowsArrJan = append(strRowsArrJan, v.Mcdd)
			}
			csvWriter.Write(strRowsArrJan)
		} else if groupType == GROUP_TYPE_MONTH {
			//Year
			strRowsArrJan := []string{}
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "年")
			for _, v := range listRange {
				strRowsArrJan = append(strRowsArrJan, v.Mcyyyy)
			}
			csvWriter.Write(strRowsArrJan)
			//Month
			strRowsArrJan = []string{}
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "月")
			for _, v := range listRange {
				strRowsArrJan = append(strRowsArrJan, v.Mcmm)
			}
			csvWriter.Write(strRowsArrJan)
		} else {
			//Year
			strRowsArrJan := []string{}
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "年")
			for _, v := range listRange {
				strRowsArrJan = append(strRowsArrJan, v.Mcyyyy)
			}
			csvWriter.Write(strRowsArrJan)
			//Week
			strRowsArrJan = []string{}
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "")
			strRowsArrJan = append(strRowsArrJan, "週")
			for _, v := range listRange {
				strRowsArrJan = append(strRowsArrJan, v.Mcweekdate)
			}
			csvWriter.Write(strRowsArrJan)
		}

		// 2. Title Row ---------------------------------------------------
		strRowsArrHeadRow := []string{}
		strRowsArrHeadRow = append(strRowsArrHeadRow, "店舗名")
		strRowsArrHeadRow = append(strRowsArrHeadRow, "入荷累計")
		strRowsArrHeadRow = append(strRowsArrHeadRow, "売上累計")
		strRowsArrHeadRow = append(strRowsArrHeadRow, "在庫数")
		strRowsArrHeadRow = append(strRowsArrHeadRow, "期間合計")
		for range listRange {
			strRowsArrHeadRow = append(strRowsArrHeadRow, "合計")
		}
		if csvWriter != nil {
			csvWriter.Write(strRowsArrHeadRow)
		}

		//3. Write data total---------------------------------------------------
		strRowsArrHeadTotal := []string{}
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, "合計")
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, strconv.Itoa(int(totalSingleItemJan.ReturnTotal)))
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, strconv.Itoa(int(totalSingleItemJan.SaleTotal)))
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, strconv.Itoa(int(totalSingleItemJan.StockCurCount)))
		strRowsArrHeadTotal = append(strRowsArrHeadTotal, strconv.Itoa(int(totalSingleItemJan.SaleTotalDate)))
		for _, v := range listRange {
			strRowsArrHeadTotal = append(strRowsArrHeadTotal, strconv.Itoa(int(totalSingleItemJan.SaleDay[totalSingleItemJan.JanCd][v.McKey])))
		}
		if csvWriter != nil {
			csvWriter.Write(strRowsArrHeadTotal)
		}

		//4. Write data detail shop ---------------------------------------------------
		for _, v := range dataSingleItemJan {
			strRowsArrRowData := []string{}
			strRowsArrRowData = append(strRowsArrRowData, v.ShopName)
			strRowsArrRowData = append(strRowsArrRowData, strconv.Itoa(int(v.ReturnTotal)))
			strRowsArrRowData = append(strRowsArrRowData, strconv.Itoa(int(v.SaleTotal)))
			strRowsArrRowData = append(strRowsArrRowData, strconv.Itoa(int(v.StockCurCount)))
			strRowsArrRowData = append(strRowsArrRowData, strconv.Itoa(int(v.SaleTotalDate)))
			for _, item := range listRange {
				strRowsArrRowData = append(strRowsArrRowData, strconv.Itoa(int(v.SaleDay[v.JanCd][item.McKey])))
			}
			if csvWriter != nil {
				csvWriter.Write(strRowsArrRowData)
			}
		}
	}
	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}
	return nil, filePath, fileName
}
