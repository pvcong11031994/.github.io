package UsersMaintenance

import (
	"WebPOS/Common"
	"WebPOS/Models/DB"
	"encoding/csv"
	"errors"
	"github.com/goframework/encode"
	"github.com/goframework/gf"
	"reflect"
	"strconv"
	"strings"
)

var (
	USER_ID            = "ユーザID"
	USER_NAME          = "ユーザ名"
	SHOP_CODE          = "店舗コード"
	SHOP_NAME          = "店舗名"
	SERVER_NAME        = "サーバー名"
	FLG_AUTH           = "権限フラグ"
	FRANCHISE_CD       = "フランチャイズコード"
	FRANCHISE_CD_GROUP = "フランチャイズグループコード"
	MENU_GROUP         = "メニュー閲覧フラグ"
	LASTEST_LOGIN_TIME = "最新ログイン時間"
	DEPT_CD            = "部署コード"
	DEPT_NAME          = "ユーザ所属部署名"
	USER_MAIL          = "ユーザメール"
	USER_TEL           = "ユーザ電話"
	USER_FAX           = "ユーザＦＡＸ"
	PASS_WORD          = "パスワード"
	FLAG_USE           = "稼動フラグ"
	CORP_CD            = "企業コード"
	CORP_NAME          = "企業名"
	REFER_ODER_TABLE   = "発注データ用テーブル名"

	HEADER_CSV_TEMP = []string{
		USER_ID,
		USER_NAME,
		SHOP_CODE,
		SHOP_NAME,
		FLG_AUTH,
		SERVER_NAME,
		FRANCHISE_CD,
		FRANCHISE_CD_GROUP,
		MENU_GROUP,
		LASTEST_LOGIN_TIME,
		DEPT_CD,
		DEPT_NAME,
		USER_MAIL,
		USER_TEL,
		USER_FAX,
		PASS_WORD,
		FLAG_USE,
		CORP_CD,
		CORP_NAME,
		REFER_ODER_TABLE,
	}
)

func CsvUploadView(ctx *gf.Context) {
	ctx.ViewData["link_action"] = PATH_MAINTENANCE_USER_REGISTER_CSV_ACTION
	ctx.View = "maintenance/user/user_register_csv.html"
}

func CsvUploadAction(ctx *gf.Context) {

	data := map[string]interface{}{}

	// Get&Read file upload
	file, err := ctx.GetUploadFile("upload-file-data")
	Common.LogErr(err)
	defer file.Close()
	fileDecoder := encode.NewDecoder("UTF-8")
	fileReader := fileDecoder.NewReader(file)
	csvReader := csv.NewReader(fileReader)
	listInsert, err := csvReader.ReadAll()
	if err != nil {
		onError(ctx, data, errors.New(Common.CSV_ERROR_0001))
		return
	}

	// Handle list error if exist
	listError := map[string]interface{}{}
	listKeyCompare := map[string]interface{}{}
	// Create data insert
	args := []interface{}{}
	for key, line := range listInsert {
		// Check format data
		if key == 0 {
			if !checkHeaderFile(line) {
				onError(ctx, data, errors.New(Common.CSV_ERROR_0001))
				return
			}
			continue
		}

		item_id := ""
		arg := []interface{}{}
		bShopCd := true
		bFranchiseCd := true
		// Create&Check data insert
		for i, item := range line {
			strKey := strconv.Itoa(key)
			if str := checkFieldRequire(i, item); str != "" {
				listError[strKey] = str
				break
			}

			if str := checkFieldAll(i, item); str != "" {
				listError[strKey] = str
				break
			}

			switch HEADER_CSV_TEMP[i] {
			case SHOP_CODE:
				if strings.TrimSpace(item) == "" {
					bShopCd = false
				}
			case FRANCHISE_CD:
				if strings.TrimSpace(item) == "" {
					bFranchiseCd = false
				}
			}
			if !bShopCd && !bFranchiseCd {
				listError[strKey] = "店舗コードまたはフランチャイズグループコードを入力してください。"
				break
			}

			// Init data insert after check
			switch HEADER_CSV_TEMP[i] {
			case USER_ID:
				listKeyCompare[item] = key
				item_id = item
			case PASS_WORD:
				item = Common.GeneratePass(item_id, strings.TrimSpace(item))
			}
			arg = append(arg, item)
		}
		if len(arg) == len(HEADER_CSV_TEMP) {
			args = append(args, arg...)
		}
	}
	// Insert data
	um := Models.UserMasterModel{ctx.DB}
	err, rs := um.CsvMultipleRegister(args)
	if err != nil {
		onError(ctx, data, errors.New(Common.CSV_ERROR_0001))
		return
	}

	totalRecordSuccess := len(listInsert) - 1 - len(listError)
	// Parse record error
	for _, errs := range rs {
		key := listKeyCompare[strings.Split(errs, "'")[1]].(int)
		strKey := strconv.Itoa(key)
		listError[strKey] = "ユーザIDが既に存在しました。"
		totalRecordSuccess--
	}
	data["listError"] = listError
	data["info"] = map[string]interface{}{
		"totalRecordSuccess": totalRecordSuccess,
	}
	data["is_success"] = "true"
	ctx.JsonResponse = data
}

func checkHeaderFile(headerLine []string) bool {
	return reflect.DeepEqual(headerLine, HEADER_CSV_TEMP)
}

func checkFieldRequire(index int, item string) string {
	msgErr := ""
	switch HEADER_CSV_TEMP[index] {
	case USER_ID:
		if strings.Contains(strings.TrimSpace(item), " ") {
			msgErr = "スペースなしでユーザIDを入力してください。"
			break
		}
		if strings.TrimSpace(item) == "" {
			msgErr = HEADER_CSV_TEMP[index] + "を入力してください。"
		}
	case USER_NAME,
		FLG_AUTH,
		MENU_GROUP,
		DEPT_CD,
		PASS_WORD:
		if strings.TrimSpace(item) == "" {
			msgErr = HEADER_CSV_TEMP[index] + "を入力してください。"
		}
		break
	}
	return msgErr
}

func checkFieldAll(index int, item string) string {
	msgErr := ""
	switch HEADER_CSV_TEMP[index] {
	case USER_ID:
		msgErr = checkFieldLength(USER_ID, item, 50)
		break
	case USER_NAME:
		msgErr = checkFieldLength(USER_NAME, item, 50)
		break
	case FLG_AUTH:
		msgErr = checkFieldLength(FLG_AUTH, item, 2)
		if msgErr == "" {
			msgErr = checkFieldNumber(FLG_AUTH, item)
		}
		break
	case SERVER_NAME:
		msgErr = checkFieldLength(SERVER_NAME, item, 20)
		break
	case SHOP_CODE:
		msgErr = checkFieldLength(SHOP_CODE, item, 10)
		break
	case SHOP_NAME:
		msgErr = checkFieldLength(SHOP_CODE, item, 20)
		break
	case FRANCHISE_CD:
		msgErr = checkFieldLength(FRANCHISE_CD, item, 10)
		break
	case FRANCHISE_CD_GROUP:
		msgErr = checkFieldLength(FRANCHISE_CD_GROUP, item, 10)
		break
	case MENU_GROUP:
		msgErr = checkFieldLength(MENU_GROUP, item, 3)
		break
	case LASTEST_LOGIN_TIME:
		msgErr = checkFieldLength(LASTEST_LOGIN_TIME, item, 14)
		if msgErr == "" {
			msgErr = checkFieldNumber(LASTEST_LOGIN_TIME, item)
		}
		break
	case DEPT_CD:
		msgErr = checkFieldLength(DEPT_CD, item, 10)
		break
	case DEPT_NAME:
		msgErr = checkFieldLength(DEPT_NAME, item, 20)
		break
	case USER_MAIL:
		msgErr = checkFieldLength(USER_MAIL, item, 50)
		break
	case USER_TEL:
		msgErr = checkFieldLength(USER_TEL, item, 15)
		if msgErr == "" {
			msgErr = checkFieldNumber(FLG_AUTH, item)
		}
		break
	case USER_FAX:
		msgErr = checkFieldLength(USER_FAX, item, 15)
		if msgErr == "" {
			msgErr = checkFieldNumber(FLG_AUTH, item)
		}
		break
	case PASS_WORD:
		msgErr = checkFieldLength(PASS_WORD, item, 50)
		break
	case FLAG_USE:
		msgErr = checkFieldLength(FLAG_USE, item, 1)
		if msgErr == "" {
			msgErr = checkFieldNumber(FLG_AUTH, item)
		}
		break
	case CORP_CD:
		msgErr = checkFieldLength(CORP_CD, item, 20)
		break
	case CORP_NAME:
		msgErr = checkFieldLength(CORP_NAME, item, 50)
		break
	case REFER_ODER_TABLE:
		msgErr = checkFieldLength(REFER_ODER_TABLE, item, 50)
		break
	}
	return msgErr
}

func checkFieldLength(itemName string, itemValue string, length int) string {
	msgErr := ""
	if len([]rune(itemValue)) > length {
		msgErr = "「" + itemName + "」は" + strconv.Itoa(length) + "文字以内入力してください。"
	}
	return msgErr
}

func checkFieldNumber(itemName string, itemValue string) string {
	msgErr := ""
	if !Common.CheckNumber(itemValue) {
		msgErr = "「" + itemName + "」は数字で入力してください。"
	}
	return msgErr
}

func onError(ctx *gf.Context, data map[string]interface{}, err error) {
	Common.LogErr(err)
	data["is_success"] = "false"
	data["message_err"] = "エラー : " + err.Error()
	ctx.JsonResponse = data
	return
}
