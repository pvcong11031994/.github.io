package RP064_BestSales_Maria

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"encoding/csv"
	"github.com/goframework/encode"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"runtime"
)

const (
	_ERROR_OVER_LIMIT_RECORD = "データが65536件を超えました。CSVでダウンロードして下さい。"
)

/*
	4. Write file ROW | COL | SUM
*/
func WriteFile4(data *RpData, searchHandleType, controlType, downloadType string) (error, string, string) {

	defer runtime.GC()
	var csvWriter *csv.Writer = nil
	var csvFile *os.File = nil

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	fileName := _REPORT_ID_DOWNLOAD + "_" + Common.CurrentDateTime()
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

	typeSearch := ""
	if controlType == CONTROL_TYPE_BOOK {
		typeSearch = "著者"
	} else if controlType == CONTROL_TYPE_MAGAZINE {
		typeSearch = "雑誌コード+月号"
	}
	// 1. Head -------------------------------------------------------
	headerRows := []string{
		"順位",
		"ＪＡＮ",
		"品名",
		typeSearch,
		"出版社名",
		"発売日",
		"本体価格",
		"入荷累計",
		"売上累計",
		"在庫数",
		"初売上日",
		"期間売上合計",
	}
	if downloadType == DOWNLOAD_TYPE_TOTAL_RESULT {
		writeNotTransition(headerRows, data, csvWriter, csvFile, controlType)
	} else if downloadType == DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION {
		writeTransition(headerRows, data, csvWriter, csvFile, controlType)
	} else if downloadType == DOWNLOAD_TYPE_TOTAL_RESULT_STORE {
		writeStore(headerRows, data, csvWriter, csvFile, controlType)
	}
	return nil, filePath, fileName
}

//Write file when フォーマット choose 集計結果
func writeNotTransition(headerRows []string, data *RpData, csvWriter *csv.Writer, csvFile *os.File, controlType string) {

	strRowsArr := []string{}
	for j := 0; j < len(headerRows); j++ {
		strRowsArr = append(strRowsArr, headerRows[j])

	}
	if csvWriter != nil {
		csvWriter.Write(strRowsArr)
	}
	// ---------------------------------------------------------------
	// 3. data
	for k, v := range data.Rows {
		if k == RPComon.REPORT_CSV_ROW_LIMIT {
			break
		}
		strData := []string{}
		for j, vItem := range v {
			if j < len(headerRows) {
				if controlType == CONTROL_TYPE_MAGAZINE && j == 3 {
					vItemFormat := Common.FormatCodeAndMonth(vItem.(string))
					strData = append(strData, vItemFormat)
				} else {
					strData = append(strData, vItem.(string))
				}
			}
		}
		if csvWriter != nil {
			csvWriter.Write(strData)
		}
	}

	// ---------------------------------------------------------------
	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}

	return
}

//Write file when フォーマット choose 集計結果+推移
func writeTransition(headerRows []string, data *RpData, csvWriter *csv.Writer, csvFile *os.File, controlType string) {

	for i, colName := range data.HeaderCols {
		strRowsArr := []string{}
		for j := 0; j < len(headerRows)-2; j++ {
			strRowsArr = append(strRowsArr, "")
		}
		strRowsArr = append(strRowsArr, colName)
		if i == 0 {
			strRowsArr = append(strRowsArr, "総合計")
		} else {
			strRowsArr = append(strRowsArr, "")
		}
		for _, v := range data.Cols[i] {
			strRowsArr = append(strRowsArr, v)
		}
		if csvWriter != nil {
			csvWriter.Write(strRowsArr)
		}
	}
	// ---------------------------------------------------------------
	// 2. Head Row ---------------------------------------------------
	strRowsArrHeadRow := headerRows
	for i := 0; i < len(data.Cols[0]); i++ {
		strRowsArrHeadRow = append(strRowsArrHeadRow, "販売数")
	}

	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadRow)
	}
	// ---------------------------------------------------------------
	// 3. data
	for k, v := range data.Rows {
		if k == RPComon.REPORT_CSV_ROW_LIMIT {
			break
		}
		strData := []string{}
		for i, vItem := range v {
			if i < len(v) {
				if controlType == CONTROL_TYPE_MAGAZINE && i == 3 {
					vItemFormat := Common.FormatCodeAndMonth(vItem.(string))
					strData = append(strData, vItemFormat)
				} else {
					strData = append(strData, vItem.(string))
				}
			}
		}
		if csvWriter != nil {
			csvWriter.Write(strData)
		}
	}
	// ---------------------------------------------------------------
	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}

	return
}

//Write file when フォーマット choose 集計結果+店舗
func writeStore(headerRows []string, data *RpData, csvWriter *csv.Writer, csvFile *os.File, controlType string) {

	//Header
	strRowsArr := []string{}
	for _, value := range headerRows {
		strRowsArr = append(strRowsArr, value)
	}
	for _, value := range data.HeaderCols {
		strRowsArr = append(strRowsArr, value)
	}
	if csvWriter != nil {
		csvWriter.Write(strRowsArr)
	}

	//Data
	for k, v := range data.Rows {
		if k == RPComon.REPORT_CSV_ROW_LIMIT {
			break
		}
		strData := []string{}
		for i, vItem := range v {
			if i < len(v) {
				if controlType == CONTROL_TYPE_MAGAZINE && i == 3 {
					vItemFormat := Common.FormatCodeAndMonth(vItem.(string))
					strData = append(strData, vItemFormat)
				} else {
					strData = append(strData, vItem.(string))
				}
			}
		}
		if csvWriter != nil {
			csvWriter.Write(strData)
		}
	}
	// ---------------------------------------------------------------
	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}

	return
}
