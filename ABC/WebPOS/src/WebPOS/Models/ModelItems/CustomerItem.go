package ModelItems

type Customer struct {
	ServerName      string `sql:"mc_server_name"`
	ShopCd          string `sql:"mc_shop_cd"`
	CustomerCatCd   string `sql:"mc_customer_cat_cd"`
	CustomerCatName string `sql:"mc_customer_cat_name"`
	GenderCd        string `sql:"mc_gender_cd"`
	GenderName      string `sql:"mc_gender_name"`
	FlgDel          string `sql:"mc_flg_del"`
}
