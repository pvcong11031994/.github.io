package RP067_ShopSales_Cloud

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"WebPOS/Models/ModelItems"
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/goframework/encode"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
)

func WriteFile1(data []SingleItem, listHeader []string, ctx *gf.Context, form QueryForm) (error, string, string) {

	defer runtime.GC()
	var csvWriter *csv.Writer = nil
	var csvFile *os.File = nil

	tmpPath, _ := filepath.Abs("./tmp")
	fileDir := filepath.Join(tmpPath, Common.CurrentDate(), Common.RandString(8))
	os.MkdirAll(fileDir, os.ModePerm)
	fileName := _REPORT_ID_DOWNLOAD + "_" + Common.CurrentDateTime()
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

	// Write Header ==================================================================================
	// Get list date header
	if form.DownloadType == DOWNLOAD_TYPE_TOTAL_RESULT_TRANSITION {
		dateSearchFrom := ""
		if form.DateFrom != "" {
			dateSearchFrom = strings.Replace(form.DateFrom, "/", "", -1)
		}
		dateSearchTo := ""
		if form.DateTo != "" {
			dateSearchTo = strings.Replace(form.DateTo, "/", "", -1)
		}
		listRange := []ModelItems.MasterCalendarItem{}

		mcmd := Models.MasterCalendarModel{ctx.DB}

		// init array header YYYY/MM/DD
		headerRangeYear := []string{}
		headerRangeMonth := []string{}
		headerRangeWeek := []string{}
		headerRangeDay := []string{}
		for i, _ := range listHeader {
			if i < len(listHeader)-2 {
				headerRangeYear = append(headerRangeYear, "")
				headerRangeMonth = append(headerRangeMonth, "")
				headerRangeWeek = append(headerRangeWeek, "")
				headerRangeDay = append(headerRangeDay, "")
			}
		}
		// Write header 推移
		switch form.GroupType {
		case GROUP_TYPE_DATE:
			headerRangeYear = append(headerRangeYear, "年", "総合計")
			headerRangeMonth = append(headerRangeMonth, "月", "")
			headerRangeDay = append(headerRangeDay, "日", "")
			listRange, err = mcmd.GetDay(dateSearchFrom, dateSearchTo)
			for _, v := range listRange {
				headerRangeYear = append(headerRangeYear, v.Mcyyyy)
				headerRangeMonth = append(headerRangeMonth, v.Mcmm)
				headerRangeDay = append(headerRangeDay, v.Mcdd)
				listHeader = append(listHeader, "販売数")
			}
			csvWriter.Write(headerRangeYear)
			csvWriter.Write(headerRangeMonth)
			csvWriter.Write(headerRangeDay)
		case GROUP_TYPE_WEEK:
			headerRangeYear = append(headerRangeYear, "年", "総合計")
			headerRangeWeek = append(headerRangeWeek, "週", "")
			listRange, err = mcmd.GetWeek(dateSearchFrom, dateSearchTo)
			for _, v := range listRange {
				headerRangeYear = append(headerRangeYear, v.Mcyyyy)
				headerRangeWeek = append(headerRangeWeek, v.Mcweekdate)
				listHeader = append(listHeader, "販売数")
			}
			csvWriter.Write(headerRangeYear)
			csvWriter.Write(headerRangeWeek)

		case GROUP_TYPE_MONTH:
			headerRangeYear = append(headerRangeYear, "年", "総合計")
			headerRangeMonth = append(headerRangeMonth, "月", "")
			listRange, err = mcmd.GetMonth(dateSearchFrom, dateSearchTo)
			for _, v := range listRange {
				headerRangeYear = append(headerRangeYear, v.Mcyyyy)
				headerRangeMonth = append(headerRangeMonth, v.Mcmm)
				listHeader = append(listHeader, "販売数")
			}
			csvWriter.Write(headerRangeYear)
			csvWriter.Write(headerRangeMonth)
		}
	}
	// Write product info
	csvWriter.Write(listHeader)

	// Write data ==================================================================================
	for _, row := range data {
		for _, item := range row.Data {
			singleRow := []string{row.ShmSharedBookStoreCode, row.ShmShopCode, row.ShmShopName}
			for i, itemDetail := range item {
				if i < 2 {
					continue
				}
				value := ""
				switch reflect.TypeOf(itemDetail).String() {
				case "string":
					value = itemDetail.(string)
				case "int64":
					valueInt := itemDetail.(int64)
					value = strconv.Itoa(int(valueInt))
				case "float64":
					valueFloat := itemDetail.(float64)
					value = strconv.Itoa(int(valueFloat))
				}
				singleRow = append(singleRow, value)
			}
			csvWriter.Write(singleRow)
		}
	}
	if csvWriter != nil {
		csvWriter.Flush()
		csvFile.Close()
	}
	return nil, filePath, fileName
}
