package ModelItems

type BQLogSearchChargingItem struct {
	CountSearchingDate int64  `sql:"count_searching_date"`
	SearchingDate      string `sql:"bqls_searching_date"`
	ServerName         string `sql:"bqls_server_name"`
	ShopCd             string `sql:"bqls_shop_cd"`
	ShopName           string `sql:"shm_shop_name"`
	FranchiseCd        string `sql:"bqls_franchise_cd"`
	FranchiseName      string `sql:"sfgm_franchise_name"`
	UserID             string `sql:"bqls_user_ID"`
	UserName           string `sql:"bqls_user_name"`
	UseTBLSize         int64  `sql:"bqls_use_TBL_size"`
	SumGCPCharging     int64  `sql:"sum_GCP_charging"`
	SumVJCharging      int64  `sql:"sum_VJ_charging"`
	Bairitsu           string `sql:"bqcs_bairitsu"`
	UserMenu           string `sql:"bqls_use_menu"`
	GCPCharging        int    `sql:"bqls_GCP_charging"`
	VJCharging         int    `sql:"bqls_VJ_charging"`
	ExecTime           int64  `sql:"bqls_exec_time"`
	AppVersion         string `sql:"bqls_app_version"`
	Handle             string `sql:"bqls_app_handle"`
	Format             string `sql:"bqls_app_format"`
	Tab                string `sql:"bqls_app_tab"`
	AppID              string `sql:"bqls_app_id"`
}

type DataChargesByUserItem struct {
	ServerName string `sql:"bqls_server_name"`

	SearchDateList   string `listCharges:"集計期間"`
	ShopCdList       string `listCharges:"店舗コード"   sql:"bqls_shop_cd_list"`
	ShopNameList     string `listCharges:"店舗名"       sql:"bqls_shop_name_list"`
	TotalSearchList  int64  `listCharges:"検索回数"     sql:"bqls_total_search_list"`
	TotalChargesList int64  `listCharges:"課金額"       sql:"bqls_total_charges_list"`

	SearchDateDetail    string `detailCharges:"集計期間"`
	ShopCdDetail        string `detailCharges:"店舗コード"    sql:"bqls_shop_cd_detail"`
	ShopNameDetail      string `detailCharges:"店舗名"        sql:"bqls_shop_name_detail"`
	SearchingDateDetail string `detailCharges:"検索日"        sql:"bqls_searching_date"`
	UserIDDetail        string `detailCharges:"ユーザID"      sql:"bqls_user_ID_detail"`
	UserNameDetail      string `detailCharges:"ユーザ名"      sql:"bqls_user_name_detail"`
	UseMenuDetail       string `detailCharges:"利用帳票名"    sql:"bqls_use_menu"`
	TotalChargesDetail  int64  `detailCharges:"課金額"        sql:"bqls_total_charges_detail"`
}
