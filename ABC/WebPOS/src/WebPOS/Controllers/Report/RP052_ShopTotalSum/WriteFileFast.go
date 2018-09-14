package RP052_ShopTotalSum

import (
	"WebPOS/Common"
	"encoding/csv"
	"github.com/goframework/encode"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

/*
	Writer file CSV
	Header: 店舗,JANコード,品名,出版社,本体価格,売上数累計,期間売上数合計,在庫数
*/
func WriteFileFast(data *DataSum) (error, string, string) {

	defer runtime.GC()
	var csvWriter *csv.Writer = nil
	var csvFile *os.File = nil
	var err error

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	fileName := _REPORT_ID + "_" + Common.CurrentDateTime()
	filePath := filepath.Join(fileDir, fileName)

	filePath = filePath + ".csv"
	fileName = fileName + ".csv"
	csvFile, err = os.Create(filePath)
	if err != nil {
		return exterror.WrapExtError(err), "", ""
	}
	csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(csvFile))
	defer csvFile.Close()
	csvWriter.UseCRLF = true

	// 1. Head -------------------------------------------------------

	strRows := "店舗,JANコード,品名,出版社,本体価格,売上数累計,期間売上数合計,在庫数"
	csvWriter.Write(strings.Split(strRows, ","))
	// ---------------------------------------------------------------
	// ---------------------------------------------------------------
	// 2. Row data
	for _, row := range data.ResultData {
		strRowsArr := []string{}
		strRowsArr = append(strRowsArr,row.ShopName)
		strRowsArr = append(strRowsArr,row.JanCd)
		strRowsArr = append(strRowsArr,row.GoodsName)
		strRowsArr = append(strRowsArr,row.PublisherName)
		strRowsArr = append(strRowsArr,strconv.FormatInt(row.Price, 10))
		strRowsArr = append(strRowsArr,strconv.FormatInt(row.SaleTotal, 10))
		strRowsArr = append(strRowsArr,strconv.FormatInt(row.SaleTotalDate, 10))
		strRowsArr = append(strRowsArr,strconv.FormatInt(row.StockCount, 10))
		csvWriter.Write(strRowsArr)
	}
	// ---------------------------------------------------------------
	csvWriter.Flush()
	csvFile.Close()

	return nil, filePath, fileName
}
