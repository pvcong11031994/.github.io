package ModelItems

type BQChargingSettingItem struct {
	ServerName    string `sql:"bqcs_server_name"`
	FranchiseCd   string
	FranchiseName string
	ShopCd        string `sql:"bqcs_shop_cd"`
	ShopName      string
	UserID        string `sql:"bqcs_user_ID"`
	UserName      string
	Bairitsu      int    `sql:"bqcs_bairitsu"`
	DeleteFlag    string `sql:"bqcs_delete_flag"`
}
