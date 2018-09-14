package Common

// Session key map

const (
	// Key generate pass secret
	SHA_PRIVATE_KEY_1 = "@#%visual%#@"
	SHA_PRIVATE_KEY_2 = "&*japan*&"

	// Common value
	HONBU_SHOP_SERVERNAME = "99999999"
	HONBU_SHOP_CD         = "99999999"
	HONBU_SHOP_NAME       = "本部"
	CHAIN_VJ              = "VJ"

	AUTH_HONBU = "1"
	AUTH_TENPO = "0"
)

//	Defined WORK_TYPE_CLASS
const (
	//=========================	業務区分	=========================//
	CATEGORY_ORDER           = "0202"
	CATEGORY_TRANSFER_EXPORT = "0235"
	CATEGORY_TRANSFER_IMPORT = "0215"
	CATEGORY_INSTRUCT_EXPORT = "0281"
	CATEGORY_INSTRUCT_IMPORT = "0282"
	CATEGORY_PROCESS_EXPORT  = "0283"
	CATEGORY_PROCESS_IMPORT  = "0284"
)

const (
	_WINDOWS_OS = "windows"
	_PATH_LOG   = "../log/WEBPOS/"
)

var listWorkType = map[string]string{
	CATEGORY_ORDER: "発注",
}

const (
	FAVORITE_MANAGEMENT_DELETE_PROCESS = "削除"
	FAVORITE_MANAGEMENT_UPDATE_PROCESS = "更新"
	FAVORITE_MANAGEMENT_LOG            = "＜%v＞が＜%v＞を＜%v＞しました。"
	INIT_MEMO                          = ""
	INIT_PRIORITY_NUMBER               = "1"
)
const (
	//ASO5402
	REPORT_MSG_ERROR = "[ERROR]"
	REPORT_MSG_RETRY = "[RETRY]"
)
