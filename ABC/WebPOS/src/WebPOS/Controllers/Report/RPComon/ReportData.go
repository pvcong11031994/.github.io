package RPComon

type ReportData struct {
	ListColKey []string
	ListRowKey []string

	Cols map[string][]string
	Rows map[string][]string
	Data map[string]map[string][]interface{}

	HeaderSum []string
	HeaderCol []string
	HeaderRow []string
	FormatSum []string

	CountResultRows int

	ShowLineFrom int
	ShowLineTo   int

	PageCount  int
	ThisPage   int
	VJCharging int
}

const REPORT_SEARCH_DATA_COLUMNS_LARGE = "列数が300を超えないように再度条件を選択してください。"

// FUJI-4677 EDIT
//const REPORT_SEARCH_RESULT_EMPTY = "該当するデータがありません。"
const REPORT_SEARCH_RESULT_EMPTY = "該当するデータがありません。条件を変更して再度検索してください。"
const REPORT_ERROR_NO_SHOP = "店舗を選択してください。"
const REPORT_ERROR_NO_SUM = "集計項目を選択してください。"
const REPORT_ERROR_NO_COL = "カラム項目を選択してください。"
const REPORT_ERROR_NO_ROW = "行項目を選択してください。"

// FUJI-4677 EDIT
//const REPORT_ERROR_DATE = "日付を再度入力してください。"
const REPORT_ERROR_DATE = "日付を入力してください。"
const REPORT_ERROR_MONTH = "月付を入力してください。"
const REPORT_ERROR_YEAR = "年を再度入力してください。 "
const REPORT_ERROR_DAY = "日を再度入力してください。"
const REPORT_ERROR_TIME = "時間を再度入力してください。"
const REPORT_ERROR_SALE_DATE = "発売日を再度入力してください。 "
const REPORT_ERROR_FIRST_SUPPLY_DATE = "初回入荷日を再度入力してください。 "
const REPORT_ERROR_FIRST_SALES_DATE = "初回販売日を再度入力してください。 "
const REPORT_ERROR_LAST_SUPPLY_DATE = "最終入荷日を再度入力してください。 "
const REPORT_ERROR_LAST_SALES_DATE = "最終販売日を再度入力してください。 "
const REPORT_ERROR_LAST_RETURN_DATE = "最終返品日を再度入力してください。 "
const REPORT_ERROR_NO_POS = "POS番号を選択してください。"
const REPORT_ERROR_WORK_TYPE = "業務区分を選択してください。 "
const REPORT_ERROR_WORK_TYPE_DISPLAY = "業務区分は必須表示です。 "
const REPORT_ERROR_LAST_SUPPLY_DATE2 = "最終仕入日を再度入力してください。 "
const REPORT_ERROR_LAST_ORDER_DATE = "最終発注日を再度入力してください。 "
const REPORT_ERROR_LAST_SALES_DATE2 = "最終売上日を再度入力してください。 "
const REPORT_ERROR_TRN_DATE = "日付を再度入力してください。 "
const REPORT_ERROR_NO_TRN_DATE = "日付を選択してください。 "
const REPORT_ERROR_NO_FLAG_SHOPCHAIN = "全国累計販売数、全国指定期間内販売数、チェーン累計販売数、チェーン指定期間内販売数を最低1つ選択してください。 "
const REPORT_ERROR_LAST_ORDER_DATE_2MONTH = "最終発注日は2か月前までの日付を選択してください。"

const REPORT_ERROR_SUPER = "初速の日を再度入力してください。"
const REPORT_ERROR_INPUT_JAN = "JANコードかISBNを入力してください。"
const REPORT_ERROR_INPUT_JAN_ARRAY_AND_JAN_SINGLE = "JANコード前方一致かJANコード完全一致を入力してください。"
const REPORT_ERROR_INPUT_JAN_SINGLE = "前方一致は6文字から検索可能です。"
const REPORT_ERROR_INPUT_MAGAZINE_CODE = "雑誌コードを入力してください。"
const REPORT_ERROR_INPUT_JAN_MAKER_CODE = "JAN/出版者記号を入力してください。"

// Type search when select btn-download-report in report
const REPORT_SEARCH_TYPE_HANDLE_CSV = "1"
const REPORT_SEARCH_TYPE_HANDLE_EXCEL = "2"

const REPORT_SEARCH_TYPE_HANDLE_TEXT = "集計実行"
const REPORT_SEARCH_TYPE_HANDLE_DOWNLOAD_TEXT = "ダウンロード"
const REPORT_SEARCH_TYPE_HANDLE_CSV_TEXT = "CSV"

const NO_KEY_FIELD = "@##NO##@"
const SUM_KEY_FIELD = "@##SUM##@"
const NO_KEY = `'` + NO_KEY_FIELD + `'`
const SUM_KEY = `'` + SUM_KEY_FIELD + `'`

const BQ_DATA_LIMIT = 10000000

const LINE_PER_PAGE = 1000

const SUM_METHOD_MASTER_BUMON int = 0
const SUM_METHOD_TRANS_BUMON int = 1

const REPORT_NAME_KEY = "_REPORT_NAME_KEY"
const REPORT_VJ_CHARGING_KEY = "_REPORT_VJ_CHARGING_KEY"

// exchange rate JPY and USD (120JPY/1USD)
const REPORT_EXCHANGE_RATE int64 = 120

// charging GOOGLE_API (5$/1TB)
const REPORT_CHARGING_GOOGLE_API_QUERY int64 = 5

const PATH_REPORT_DOWN_LOAD_LINK = "/report/download"

const REPORT_QUERY_JOB_ID = "_REPORT_QUERY_JOB_ID"
const REPORT_QUERY_JOB_ID_COUNT = "_REPORT_QUERY_JOB_ID_COUNT"

// システムエラー
const REPORT_ERROR_SYSTEM = "システムエラーです。しばらくしてから再度操作をお願いします。"
const REPORT_ERROR_PATH_HTML = "/report/error/error.html"
const REPORT_ERROR_SYSTEM_VIEW = "err_msg"

// CSV ROW LIMIT
const REPORT_CSV_ROW_LIMIT = 50000

//共通仕様 ※バリデーションチェック
//期間 (日別)
const REPORT_ERROR_DATE_FORMAT = "期間は「yyyy/mm/dd」で入力してください。"
const REPORT_ERROR_100_DATE = "日別の期間は100日以内に変更してください。"

//期間 (週別)
const REPORT_ERROR_WEEK_FORMAT = "期間は「yyyy/mm/dd～dd」で入力してください。"
const REPORT_ERROR_100_WEEK = "週別の期間は100週以内に変更してください。"

//期間 (月別)
const REPORT_ERROR_MONTH_FORMAT = "期間は「yyyy/mm」で入力してください。"
const REPORT_ERROR_100_MONTH = "月別の期間は100ヶ月以内に変更してください。"

//出版社コード 		| ※出版社コードが必須の画面のみ。
const REPORT_ERROR_MAKE_CODE_BLANK = "出版社コードを入力してください。"

//出版者記号		| ※出版者記号が必須の画面のみ。
const REPORT_ERROR_JAN_MAKE_CODE_BLANK = "出版者記号を入力してください。"

//雑誌コード		| ※雑誌コードが必須の画面のみ。
const REPORT_ERROR_MAGAZINE_CODE_BLANK = "	雑誌コードを入力してください。"

// PG_ASO-5398 [BA]mBAWEB-v02a 【至急】単品推移の期間の制限を変更
// Common Limit search date
const REPORT_LIMIT_DATE_SEARCH = 100
const REPORT_ERROR_LIMIT_DATE = "日別の期間は100日以内に変更してください。"
// Common Limit search week
const REPORT_LIMIT_WEEK_SEARCH = 30
const REPORT_ERROR_LIMIT_WEEK = "週別の期間は30週以内に変更してください。"
// Common Limit search Month
const REPORT_LIMIT_MONTH_SEARCH = 13
const REPORT_ERROR_LIMIT_MONTH = "月別の期間は13ヶ月以内に変更してください。"
// Len 出版社コード
const REPORT_ERROR_LEN_MAKER_CODE = "出版社コードは4桁の数字で入力してください。"
// Select item キーワード
const REPORT_ERROR_SELECTED  = "アイテムは和書か雑誌を選択してください。"
//お気に入り管理
const REPORT_INSERT_FAVORITE_SUCCESS = "お気に入り商品を新規登録しました。"
const REPORT_UPDATE_FAVORITE_SUCCESS = "お気に入り商品を保存しました。"
const REPORT_DELETE_FAVORITE_SUCCESS = "お気に入り商品を削除しました。"
const REPORT_INSERT_FAVORITE_FAIL = "お気に入り商品の登録に失敗しました。"
const REPORT_UPDATE_FAVORITE_FAIL = "お気に入り商品の更新に失敗しました。"
const REPORT_DELETE_FAVORITE_FAIL = "お気に入り商品の削除に失敗しました。"
//商品検索
const REPORT_ERROR_LEN_KEY_WORD  = "キーワードは2文字以上入力してください。"
// Common Limit search shop
const REPORT_LIMIT_SHOP_SEARCH = 100
const REPORT_ERROR_LIMIT_SHOP_SEARCH = "Overload data"
// Download file size limit
//const REPORT_DOWNLOAD_FILE_SIZE_OVER = "File size: %sMB. It's over size."
const REPORT_DOWNLOAD_FILE_SIZE_OVER = `ダウンロード件数：%s件。
データ量が多過ぎて、全件ダウンロードすることができません。
条件を変更してください。`