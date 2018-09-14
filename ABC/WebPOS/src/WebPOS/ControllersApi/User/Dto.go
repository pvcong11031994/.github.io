package User

import (
	"WebPOS/ControllersApi/Utils"
	"WebPOS/ControllersApi/Utils/verify"
	"encoding/json"
	"errors"
)

const (
	SEARCH     = "0"
	REGISTER   = "1"
	UPDATE     = "2"
	DELETE     = "3"
	RESET_PASS = "4"
)

const (
	PROCESS_NORMAL       = "0"
	PROCESS_ERROR_APP    = "1"
	PROCESS_ERROR_SYSTEM = "99"
)

const (
	// Const info
	MSG_INFO_USER_EXISTS     = "指定されたユーザIDは既に登録済みです。"
	MSG_INFO_USER_NOT_EXISTS = "指定されたユーザIDは存在しません。"
	MSG_INFO_USER_SEARCH     = "ユーザ検索処理が正常に終了しました。"
	MSG_INFO_USER_REGISTER   = "ユーザ登録処理が完了しました。"
	MSG_INFO_USER_UPDATE     = "ユーザ更新処理が完了しました。"
	MSG_INFO_USER_DELETE     = "ユーザ削除処理が完了しました。"
	MSG_INFO_USER_RESET_PASS = "ユーザパスワードリセットが完了しました。"
	/*-----------------------------------------------------------------------------------------------------------------------------*/

	// Const error
	MSG_ERR_INTERNAL_SERVER  = "インターナルサーバエラー。"
	MSG_ERR_ROLE_NOT_GRANTED = "APIKEYが一致しません。"
	/*-----------------------------------------------------------------------------------------------------------------------------*/
)

var (
	messageResponse = map[string]string{
		SEARCH:     MSG_INFO_USER_SEARCH,
		REGISTER:   MSG_INFO_USER_REGISTER,
		UPDATE:     MSG_INFO_USER_UPDATE,
		DELETE:     MSG_INFO_USER_DELETE,
		RESET_PASS: MSG_INFO_USER_RESET_PASS,
	}
)

type (
	RequestDto struct {
		ApiKey           string `json:"um_api_key"             form:"um_api_key"                validate:"nonzero,max=64"                                                           message:"APIKEYが不正です。"`
		Exec             string `json:"um_exec"                form:"um_exec"                   validate:"nonzero,regexp=^(0|1|2|3|4)$"                                             message:"処理区分が不正です。"`
		PageNum          string `json:"um_pageNum"             form:"um_pageNum"                validate:"min=1,max=8,regexp=^[0-9]+$"                      tagCase:"0"             message:"検索ページ番号が不正です。"`
		DisplayCount     string `json:"um_displayCount"        form:"um_displayCount"           validate:"min=1,max=5,regexp=^[0-9]+$"                      tagCase:"0"             message:"検索表示件数が不正です。"`
		UserID           string `json:"um_user_ID"             form:"um_user_ID"                validate:"nonzero" require:"max=50,regexp=^[a-zA-Z0-9_]+$"  tagCase:"1,2,3,4"       message:"ユーザIDが不正です。"`
		UserName         string `json:"um_user_name"           form:"um_user_name"              validate:"nonzero"                                          tagCase:"1,2"           message:"ユーザ名が不正です。"`
		FlgAuth          string `json:"um_flg_auth"            form:"um_flg_auth"               validate:"nonzero" require:"regexp=^(0|1)$"                 tagCase:"1,2"           message:"権限フラグが不正です。"`
		ShopChainCd      string `json:"um_shop_chain_cd"       form:"um_shop_chain_cd"`
		FranchiseCd      string `json:"um_franchise_cd"        form:"um_franchise_cd"`
		FranchiseGroupCd string `json:"um_franchise_group_cd"  form:"um_franchise_group_cd"`
		ServerName       string `json:"um_server_name"         form:"um_server_name"`
		ShopCd           string `json:"um_shop_cd"             form:"um_shop_cd"                validate:"nonzero" require:"max=10,regexp=^[a-zA-Z0-9]+$"  tagCase:"1,2"            message:"店舗コードが不正です。"`
		ShopName         string `json:"um_shop_name"           form:"um_shop_name"`
		FlgMenuGroup     string `json:"um_flg_menu_group"      form:"um_flg_menu_group"         validate:"nonzero" require:"max=3,regexp=^[a-zA-Z0-9]+$"   tagCase:"1,2"            message:"メニュー閲覧フラグが不正です。"`
		DeptCd           string `json:"um_dept_cd"             form:"um_dept_cd"                validate:"nonzero" require:"max=10,regexp=^[a-zA-Z0-9_]+$" tagCase:"1,2"            message:"部署コードが不正です。"`
		DeptName         string `json:"um_dept_name"           form:"um_dept_name"`
		UserMail         string `json:"um_user_mail"           form:"um_user_mail"`
		UserPhone        string `json:"um_user_phone"          form:"um_user_phone"`
		UserXerox        string `json:"um_user_xerox"          form:"um_user_xerox"`
		FlgUse           string `json:"um_flg_use"             form:"um_flg_use"                validate:"nonzero" require:"regexp=^(0|1)$"                tagCase:"1,2"            message:"稼動フラグが不正です。"`
		CorpCd           string `json:"um_corp_cd"             form:"um_corp_cd"`
		CorpName         string `json:"um_corp_name"           form:"um_corp_name"`
		ReferOrderTable  string `json:"um_refer_order_table"   form:"um_refer_order_table"`
		PublisherFlg     string `json:"um_publisher_flg"       form:"um_publisher_flg"`
		PublisherId      string `json:"um_publisher_id"        form:"um_publisher_id"`
		OptionFlgReturn  string `json:"um_option_flg_return"   form:"um_option_flg_return"`
	}

	ResponseDto struct {
		ResultCode    string `json:"resultCode,require"`
		ResultMessage string `json:"resultMessage,require"`

		// 処理区分：「1ー4」
		InitPw string `json:"initPw,omitempty"`

		// 処理区分：「0」
		// 検索結果
		TotalCount       string      `json:"totalCount,omitempty"`
		SearchResultList []ResultDto `json:"searchResultList,omitempty"`
		PageNum          string      `json:"pageNum,omitempty"`
		DisplayCount     string      `json:"displayCount,omitempty"`
	}

	ResultDto struct {
		CreateDate       string `json:"um_create_date"`
		UpdateDate       string `json:"um_update_date"`
		UserID           string `json:"um_user_ID"`
		UserName         string `json:"um_user_name"`
		FlgAuth          string `json:"um_flg_auth"`
		ShopChainCd      string `json:"um_shop_chain_cd"`
		FranchiseCd      string `json:"um_franchise_cd"`
		FranchiseGroupCd string `json:"um_franchise_group_cd"`
		ServerName       string `json:"um_server_name"`
		ShopCd           string `json:"um_shop_cd"`
		ShopName         string `json:"um_shop_name"`
		FlgMenuGroup     string `json:"um_flg_menu_group"`
		LatestLoginTime  string `json:"um_latest_login_time"`
		DeptCd           string `json:"um_dept_cd"`
		DeptName         string `json:"um_dept_name"`
		UserMail         string `json:"um_user_mail"`
		UserPhone        string `json:"um_user_phone"`
		UserXerox        string `json:"um_user_xerox"`
		FlgUse           string `json:"um_flg_use"`
		CorpCd           string `json:"um_corp_cd"`
		CorpName         string `json:"um_corp_name"`
		ReferOrderTable  string `json:"um_refer_order_table"`
		PublisherFlg     string `json:"um_publisher_flg"`
		PublisherId      string `json:"um_publisher_id"`
		PwExpiringDate   string `json:"um_pw_expiring_date"`
		PwChangeFlg      string `json:"um_pw_change_flg"`
		OptionFlgReturn  string `json:"um_option_flg_return"`
	}
)

// Verify RequestDto with rule
func (this *RequestDto) VerifyInput() error {
	validator := verify.NewValidator()
	validator.SetCase(this.Exec)
	if errs := validator.Validate(this); errs != nil {
		errorStr := validator.ParseErr(errs.Error())
		return errors.New(errorStr)
	}
	return nil
}

// Generate ResponseDto with status code and exec flag
func (this *ResponseDto) CreateResponse(exec string) []byte {
	//
	responseObj := this.clone(exec)
	responseByte, _ := json.Marshal(responseObj)
	//
	return ApiUtils.ResponseWithCallbackEmbed("callback", responseByte)
}

// Make last response with exec flag
func (this *ResponseDto) clone(exec string) ResponseDto {
	//
	cloneObj := ResponseDto{}
	cloneObj.ResultCode = this.ResultCode
	cloneObj.ResultMessage = this.ResultMessage
	//
	if this.ResultCode != PROCESS_NORMAL {
		return cloneObj
	}
	//
	cloneObj.ResultMessage = messageResponse[exec]
	switch exec {
	// 処理区分：「0」
	case SEARCH:
		cloneObj.TotalCount = this.TotalCount
		cloneObj.SearchResultList = this.SearchResultList
		cloneObj.PageNum = this.PageNum
		cloneObj.DisplayCount = this.DisplayCount
	// 処理区分：「1ー4」
	case REGISTER, RESET_PASS:
		cloneObj.InitPw = this.InitPw
	}
	return cloneObj
}
