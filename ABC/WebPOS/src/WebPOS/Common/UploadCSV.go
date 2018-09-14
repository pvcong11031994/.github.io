package Common

import (
	"encoding/csv"
	"fmt"
	"github.com/goframework/encode"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"reflect"
	"strconv"
	"strings"
)

type CsvFileUpload struct {
	AppendDate     bool
	FileName       string
	DbFieldInsert  []string
	Header         []string
	ItemFormat     map[string]Item
	DecoderFormat  string
	FuncRegist     func([]string, []interface{}) (error, []string)
	FuncFormatItem func(int, *[]string, *string, []string)
}

const (
	TYPE_STRING   = 0
	TYPE_INTERGER = 1
)

type Item struct {
	DB_Type   int
	DB_Length int
	DB_Null   bool
	DB_Blank  bool
}

//Defined ErrorMessage
const (
	ERROR_FORMAT_FILE     = "ファイル形式が不正です。"
	ERROR_DUPLICATE       = "既に同じデータがあります。"
	ERROR_INTERNAL_SERVER = "インターナルサーバエラー"
)

func (this *CsvFileUpload) checkHeaderFile(headerLine []string) bool {
	return reflect.DeepEqual(headerLine, this.Header)
}

func setError(ctx *gf.Context, strErr string) map[string]interface{} {
	strErr = IF(strings.TrimSpace(strErr) == "", ERROR_INTERNAL_SERVER, strErr).(string)
	return map[string]interface{}{
		"success": false,
		"msg_err": "エラー : " + strErr,
	}
}

func (this *CsvFileUpload) InsertFileUploadData(ctx *gf.Context) map[string]interface{} {
	// Get&Read file upload
	file, err := ctx.GetUploadFile(this.FileName)
	LogErr(exterror.WrapExtError(err))
	defer file.Close()
	fileDecoder := encode.NewDecoder(this.DecoderFormat)
	fileReader := fileDecoder.NewReader(file)
	csvReader := csv.NewReader(fileReader)
	listInsert, err := csvReader.ReadAll()
	LogErr(exterror.WrapExtError(err))
	if err != nil {
		return setError(ctx, ERROR_FORMAT_FILE)
	}

	// Handle list error if exist
	listError := map[string]interface{}{}
	listKeyCompare := map[string]interface{}{}
	// Create data insert
	args := []interface{}{}
	date := CurrentDate()
	for key, line := range listInsert {
		// Check format data
		if key == 0 {
			if !this.checkHeaderFile(line) {
				return setError(ctx, ERROR_FORMAT_FILE)
			}
			continue
		}

		arrItemsKey := []string{}
		arg := []interface{}{}
		if this.AppendDate {
			arg = append(arg, date, date)
		}
		// Create&Check data insert
		strKey := ""
		for i, item := range line {
			strKey = strconv.Itoa(key)
			if str := this.checkItem(i, item); str != "" {
				listError[strKey] = str
				break
			}
			this.FuncFormatItem(i, &arrItemsKey, &item, line)
			arg = append(arg, item)
		}
		lenArg := len(arg)
		if this.AppendDate {
			lenArg -= 2
		}
		strKeyDB := strings.Join(arrItemsKey, "")
		if strKeyDB == "" && lenArg == len(this.Header) {
			listError[strKey] = "PRIMARY_KEYを入力してください。"
		}
		if strKeyDB != "" && lenArg == len(this.Header) {
			args = append(args, arg...)
		}
		strKeyDB = strings.Join(arrItemsKey, "-")
		this.createKeyCompare(&strKeyDB)
		listKeyCompare[strKeyDB] = key
	}

	// Insert data
	err, rs := this.FuncRegist(this.DbFieldInsert, args)
	LogErr(err)
	if err != nil {
		return setError(ctx, "データ挿入にエラーが発生しました。")
	}

	totalRecordSuccess := len(listInsert) - 1 - len(listError)
	// Parse record error
	for _, errs := range rs {
		strKeyDB := strings.Split(errs, "'")[1]
		temp := listKeyCompare[strKeyDB]
		key := temp.(int)
		strKey := strconv.Itoa(key)
		listError[strKey] = ERROR_DUPLICATE
		totalRecordSuccess--
	}
	info := map[string]interface{}{
		"totalRecordSuccess": totalRecordSuccess,
	}
	return map[string]interface{}{
		"success":   true,
		"info":      info,
		"listError": listError,
	}
}

func (this *CsvFileUpload) checkItem(index int, value string) string {
	if len(this.ItemFormat) == 0 {
		return ""
	}
	item := this.ItemFormat[this.Header[index]]
	if item == (Item{}) {
		return ""
	}
	// Check NULL
	if !item.DB_Null && !item.DB_Blank && strings.TrimSpace(value) == "" {
		return this.Header[index] + "を入力してください。"
	}
	// Check LENGTH
	if item.DB_Length > 0 && len(value) > item.DB_Length {
		return fmt.Sprintf("%vは%v桁を入力してください。", this.Header[index], item.DB_Length)
	}
	// Check Type
	return ""
}

func (this *CsvFileUpload) createKeyCompare(str *string) {
	if len(*str) >= 64 {
		decoder := encode.NewDecoder(this.DecoderFormat)
		num := 64
		strTemp := *str
		for {
			temp, flg := decoder.ConvertStringOK(strTemp[:num])
			if flg {
				*str = temp
				return
			}
			num--
		}
	}
}
