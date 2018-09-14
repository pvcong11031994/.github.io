package ModelItems

import (
	"encoding/gob"
)

type UserItem struct {
	//
	CreateDate       string `sql:"um_create_date"         json:"um_create_date"`
	UpdateDate       string `sql:"um_update_date"         json:"um_update_date"`
	UserID           string `sql:"um_user_ID"             json:"um_user_ID"`
	UserName         string `sql:"um_user_name"           json:"um_user_name"`
	FlgAuth          string `sql:"um_flg_auth"            json:"um_flg_auth"`
	ShopChainCd      string `sql:"um_shop_chain_cd"       json:"um_shop_chain_cd"`
	FranchiseCd      string `sql:"um_franchise_cd"        json:"um_franchise_cd"`
	FranchiseGroupCd string `sql:"um_franchise_group_cd"  json:"um_franchise_group_cd"`
	ServerName       string `sql:"um_server_name"         json:"um_server_name"`
	ShopCd           string `sql:"um_shop_cd"             json:"um_shop_cd"`
	ShopName         string `sql:"um_shop_name"           json:"um_shop_name"`
	FlgMenuGroup     string `sql:"um_flg_menu_group"      json:"um_flg_menu_group"`
	LatestLoginTime  string `sql:"um_latest_login_time"   json:"um_latest_login_time"`
	DeptCd           string `sql:"um_dept_cd"             json:"um_dept_cd"`
	DeptName         string `sql:"um_dept_name"           json:"um_dept_name"`
	UserMail         string `sql:"um_user_mail"           json:"um_user_mail"`
	UserPhone        string `sql:"um_user_phone"          json:"um_user_phone"`
	UserXerox        string `sql:"um_user_xerox"          json:"um_user_xerox"`
	UserPass         string `sql:"um_user_pass"           json:"_no_return"`
	FlgUse           string `sql:"um_flg_use"             json:"um_flg_use"`
	CorpCd           string `sql:"um_corp_cd"             json:"um_corp_cd"`
	CorpName         string `sql:"um_corp_name"           json:"um_corp_name"`
	ReferOrderTable  string `sql:"um_refer_order_table"   json:"um_refer_order_table"`
	PublisherFlg     string `sql:"um_publisher_flg"       json:"um_publisher_flg"`
	PublisherId      string `sql:"um_publisher_id"        json:"um_publisher_id"`
	PwExpiringDate   string `sql:"um_pw_expiring_date"    json:"um_pw_expiring_date"`
	PwChangeFlg      string `sql:"um_pw_change_flg"       json:"um_pw_change_flg"`
	OptionFlgReturn  string `sql:"um_option_flg_return"   json:"um_option_flg_return"`

	Design *ControlReportDesignItem
}

type UserRoleItem struct {
	UserID           string `sql:"um_user_ID"`
	ShopCD           string `sql:"um_shop_cd"`
	ServerName       string `sql:"um_server_name"`
	IsHonbu          bool   `sql:"um_flg_auth"`
	ChainCD          string `sql:"um_shop_chain_cd"`
	ListChainCD      []string
	FranchiseCd      string `sql:"um_franchise_cd"`
	FranchiseGroupCd string `sql:"um_franchise_group_cd"`
}

func init() {
	gob.Register(UserItem{})
}
