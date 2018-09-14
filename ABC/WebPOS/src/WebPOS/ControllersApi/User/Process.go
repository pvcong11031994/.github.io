package User

import (
	"WebPOS/Common"
	"WebPOS/ControllersApi/Utils"
	"WebPOS/Models/DB"
	"encoding/json"
	"errors"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strconv"
)

// Handle action for route : /RspAPI/user
func doProcess(ctx *gf.Context) {
	//
	responseDto := ResponseDto{}
	responseDto.ResultCode = PROCESS_ERROR_APP
	//
	formInput := RequestDto{}
	ctx.Form.ReadStruct(&formInput)

	// APIKEY認証・処理区分 + 入力チェック
	if err := formInput.VerifyInput(); err != nil {
		//
		responseDto.ResultMessage = err.Error()
		ctx.JsonPResponse = responseDto.CreateResponse(formInput.Exec)
		return
	}
	// Verify Api key
	if ApiUtils.VerifyApiKey(formInput.ApiKey) == false {
		//
		responseDto.ResultMessage = MSG_ERR_ROLE_NOT_GRANTED
		ctx.JsonPResponse = responseDto.CreateResponse(formInput.Exec)
		return
	}

	//
	umModel := Models.NewUserMaster(ctx.DB)
	userExists := false
	if formInput.Exec != SEARCH {
		var errExists error
		// Check exists ユーザID
		userExists, errExists = umModel.IsExist(formInput.UserID)
		if errExists != nil {
			renderErrorSystem(ctx, &responseDto, errExists, formInput.Exec)
			return
		}
		responseDto.ResultMessage = MSG_INFO_USER_NOT_EXISTS
		if userExists {
			responseDto.ResultMessage = MSG_INFO_USER_EXISTS
		}
	}

	// プロセス
	var err error
	responseDto.ResultCode = PROCESS_NORMAL
	switch formInput.Exec {
	// 処理区分：「0」
	case SEARCH:
		err = search(formInput, &responseDto, umModel)

	// 処理区分：「1」
	case REGISTER:
		if userExists {
			responseDto.ResultCode = PROCESS_ERROR_APP
			break
		}
		var parameter []byte
		parameter, err = json.Marshal(formInput)
		if err != nil {
			break
		}
		responseDto.InitPw, err = umModel.AddUser(parameter)

	// 処理区分：「2」
	case UPDATE:
		if !userExists {
			responseDto.ResultCode = PROCESS_ERROR_APP
			break
		}
		var parameter []byte
		parameter, err = json.Marshal(formInput)
		if err != nil {
			break
		}
		err = umModel.UpdateUser(parameter)

	// 処理区分：「3」
	case DELETE:
		if !userExists {
			responseDto.ResultCode = PROCESS_ERROR_APP
			break
		}
		err = umModel.DeleteUser(formInput.UserID)

	// 処理区分：「4」
	case RESET_PASS:
		if !userExists {
			responseDto.ResultCode = PROCESS_ERROR_APP
			break
		}
		responseDto.InitPw, err = umModel.ResetPassword(formInput.UserID)

	default:
		err = exterror.WrapExtError(errors.New(MSG_ERR_INTERNAL_SERVER))
	}

	// アウトプット
	if err != nil {
		renderErrorSystem(ctx, &responseDto, err, formInput.Exec)
		return
	}
	// Success
	ctx.JsonPResponse = responseDto.CreateResponse(formInput.Exec)
}

// 処理区分 : 検索
func search(formInput RequestDto, responseDto *ResponseDto, umModel *Models.UserMasterAPI_Model) error {
	//
	offsetCount, limitCount := ApiUtils.ParseOffset(formInput.PageNum, formInput.DisplayCount)
	parameter := map[string]interface{}{
		"um_user_id":   formInput.UserID,
		"um_user_name": formInput.UserName,
		"um_dept_cd":   formInput.DeptCd,
		"um_flg_use":   formInput.FlgUse,
		"LIMIT":        limitCount,
		"OFFSET":       offsetCount,
	}
	listUser, totalCount, err := umModel.GetUsers(parameter)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	byteRaw, _ := json.Marshal(*listUser)
	err = json.Unmarshal(byteRaw, &responseDto.SearchResultList)
	//
	if err == nil {
		responseDto.TotalCount = strconv.Itoa(totalCount)
		responseDto.PageNum = formInput.PageNum
		responseDto.DisplayCount = formInput.DisplayCount
	}
	//
	return exterror.WrapExtError(err)
}

// Make error system
func renderErrorSystem(ctx *gf.Context, responseDto *ResponseDto, err error, flag string) {
	Common.LogErr(err)
	responseDto.ResultCode = PROCESS_ERROR_SYSTEM
	responseDto.ResultMessage = MSG_ERR_INTERNAL_SERVER
	ctx.JsonPResponse = responseDto.CreateResponse(flag)
}
