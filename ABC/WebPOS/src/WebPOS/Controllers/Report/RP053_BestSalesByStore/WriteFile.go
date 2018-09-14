package RP053_BestSalesByStore

import (
	"WebPOS/Common"
	"WebPOS/Controllers/Report/RPComon"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/goframework/encode"
	"github.com/goframework/gf/exterror"
	"os"
	"path/filepath"
	"runtime"
)

const (
	_EXCEL_ORDER_DATA_SHEET_NAME = "データ"
	_EXCEL_RECORD_LIMIT          = 65536
	_EXCEL_COLUMN_LIMIT          = 256
	_ERROR_OVER_LIMIT_RECORD     = "データが65536件を超えました。CSVでダウンロードして下さい。"
	_ERROR_OVER_LIMIT_COLUMN     = "データが256カラムを超えました。CSVでダウンロードして下さい。"
)

/*
	4. Write file ROW | COL | SUM
*/
func WriteFile4(data *RPComon.ReportData, searchHandleType, groupType string) (error, string, string) {
	defer runtime.GC()
	var csvWriter *csv.Writer = nil
	var excelWriter *Common.SimpleExcelFile = nil
	var csvFile *os.File = nil
	var err error

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
	} else {
		filePath = filePath + ".xlsx"
		fileName = fileName + ".xlsx"
		excelWriter, err = Common.NewSimpleExcelFile(filePath, _EXCEL_ORDER_DATA_SHEET_NAME)
		defer excelWriter.Close()
		if err != nil {
			return exterror.WrapExtError(err), "", ""
		}
	}

	// 1. Head -------------------------------------------------------
	for i, colName := range data.HeaderCol {
		strRowsArr := []string{}
		for j, _ := range data.HeaderRow {
			if j < len(data.HeaderRow)-1 {
				strRowsArr = append(strRowsArr, "")
			}
		}
		strRowsArr = append(strRowsArr, colName)
		for k, _ := range data.HeaderSum {
			if i == 0 && k == 0 {
				strRowsArr = append(strRowsArr, "販売期間合計")
			} else {
				strRowsArr = append(strRowsArr, "")
			}
		}
		for _, colKey := range data.ListColKey {
			for k, _ := range data.HeaderSum {
				if k == 0 {
					strRowsArr = append(strRowsArr, data.Cols[colKey][i] + "日目")
				} else {
					strRowsArr = append(strRowsArr, "")
				}
			}
		}
		if csvWriter != nil {
			csvWriter.Write(strRowsArr)
		}
		if excelWriter != nil {
			arrColumnHeader := Common.ToInterfaceArray(strRowsArr)
			if len(arrColumnHeader) > _EXCEL_COLUMN_LIMIT {
				excelWriter.Destroy()
				return errors.New(_ERROR_OVER_LIMIT_COLUMN), "", ""
			}
			excelWriter.WriteData(arrColumnHeader)
		}
	}
	// ---------------------------------------------------------------
	// 2. Head Row ---------------------------------------------------
	strRowsArrHeadRow := []string{}
	for _, rowsName := range data.HeaderRow {
		strRowsArrHeadRow = append(strRowsArrHeadRow, rowsName)
	}
	for _, sumName := range data.HeaderSum {
		strRowsArrHeadRow = append(strRowsArrHeadRow, sumName)
	}
	for range data.ListColKey {
		for _, sumName := range data.HeaderSum {
			strRowsArrHeadRow = append(strRowsArrHeadRow, sumName)
		}
	}
	if csvWriter != nil {
		csvWriter.Write(strRowsArrHeadRow)
	}
	if excelWriter != nil {
		excelWriter.WriteData(Common.ToInterfaceArray(strRowsArrHeadRow))
	}
	// ---------------------------------------------------------------
	// 4. Row data
	countRecord := 0
	for _, rowKey := range data.ListRowKey {
		strRowsArr := []string{}
		for _, rows := range data.Rows[rowKey] {
			strRowsArr = append(strRowsArr, rows)
		}
		for _, val := range data.Data[rowKey][RPComon.SUM_KEY_FIELD] {
			strRowsArr = append(strRowsArr, fmt.Sprintf("%v", val))
		}
		for _, colKey := range data.ListColKey {
			if data.Data[rowKey][colKey] == nil {
				for range data.HeaderSum {
					strRowsArr = append(strRowsArr, "")
				}
			} else {
				for _, data := range data.Data[rowKey][colKey] {
					strRowsArr = append(strRowsArr, fmt.Sprintf("%v", data))
				}
			}
		}
		if csvWriter != nil {
			csvWriter.Write(strRowsArr)
		}
		if excelWriter != nil {
			if countRecord >= _EXCEL_RECORD_LIMIT {
				excelWriter.Destroy()
				return errors.New(_ERROR_OVER_LIMIT_RECORD), "", ""
			}
			excelWriter.WriteData(Common.ToInterfaceArray(strRowsArr))
		}
		countRecord = countRecord + 1
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
