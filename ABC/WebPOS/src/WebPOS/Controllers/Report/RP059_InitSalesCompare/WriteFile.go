package RP059_InitSalesCompare

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"encoding/csv"
	"github.com/goframework/encode"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
)

/*
	Write file ROW | COL | SUM
*/
func WriteFileCSV(data *RpData, controlType string, ctx *gf.Context, form QueryForm) (error, string, string) {

	defer runtime.GC()
	var csvWriter *csv.Writer = nil
	var csvFile *os.File = nil

	// process create path (file csv) in server
	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	fileName := _REPORT_ID + "_" + Common.CurrentDateTime()
	filePath := filepath.Join(fileDir, fileName)
	filePath = filePath + ".csv"
	fileName = fileName + ".csv"
	csvFile, err := os.Create(filePath)
	if err != nil {
		return exterror.WrapExtError(err), "", ""
	}
	csvWriter = csv.NewWriter(encode.NewEncoder(encode.ENCODER_SHIFT_JIS).NewWriter(csvFile))
	defer csvFile.Close()
	csvWriter.UseCRLF = true

	// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
	//// 1. Main Data -------------------------------------------------------
	//if controlType == CONTROL_TYPE_MAGAZINE {
	//	csvWriter.Write([]string{"雑誌名", data.MagazineName})
	//	csvWriter.Write([]string{"出版社名", data.MakerName})
	//	// write blank row
	//	csvWriter.Write([]string{})
	//}
	//listJanKey := writeDataIntoCSV(data, csvWriter)
	//// 2. Detail Data---------------------------------------------------------------
	//sql, listHeader := buildDetailForCSV(form, ctx, listJanKey)
	//err = queryDataDetailForCSV(ctx, sql, csvWriter, listHeader, form)
	// ---------------------------------------------------------------
	// 1. Main Data -------------------------------------------------------
	listJanKey := writeDataIntoCSV(data, csvWriter, controlType)
	sql, _, listHeader := buildDetailForCSV(form, ctx, listJanKey)
	err = queryDataDetailForCSV(ctx, sql, csvWriter, listJanKey, listHeader, form)
	// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
	csvWriter.Flush()
	return nil, filePath, fileName
}

//Write file
// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
//func writeDataIntoCSV(data *RpData, csvWriter *csv.Writer) []string {
func writeDataIntoCSV(data *RpData, csvWriter *csv.Writer, controlType string) []string {
	// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
	csvWriter.Write(data.HeaderCols)
	// ---------------------------------------------------------------
	// 2. data
	listJanKey := []string{}
	for rowIndex, rowValue := range data.Rows {
		if rowIndex == RPComon.REPORT_CSV_ROW_LIMIT {
			break
		}

		strData := []string{}
		if rowIndex < 15 {
			strData = append(strData, ListRank[rowIndex])
		} else {
			strData = append(strData, ListRank[15])
		}

		for keyItem, vItem := range rowValue {
			value := ""
			switch reflect.TypeOf(vItem).String() {
			case "string":
				value = vItem.(string)
			case "int64":
				valueInt64 := vItem.(int64)
				value = strconv.Itoa(int(valueInt64))
			case "float64":
				valueFloat := vItem.(float64)
				value = strconv.Itoa(int(valueFloat))
			case "int":
				valueInt := vItem.(int)
				value = strconv.Itoa(valueInt)
			}
			// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT START
			//strData = append(strData, value)
			if controlType == CONTROL_TYPE_MAGAZINE && keyItem == 2 {
				vItemFormat := Common.FormatCodeAndMonth(vItem.(string))
				strData = append(strData, vItemFormat)
			} else {
				strData = append(strData, value)
			}
			// ASO-5598 [BA]mBAWEB-v11a 初速比較：CLOUD化 - EDIT END
			// get JAN (index = 1)
			if keyItem == 0 {
				listJanKey = append(listJanKey, value)
			}
		}
		csvWriter.Write(strData)
	}
	return listJanKey
}
