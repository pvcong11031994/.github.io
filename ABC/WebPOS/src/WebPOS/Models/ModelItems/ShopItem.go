package ModelItems

type ShopItem struct {
	ShopCD              string `sql:"shm_shop_cd"`
	ShopName            string `sql:"shm_shop_name"`
	ServerName          string `sql:"shm_server_name"`
	ShopChainCD         string `sql:"shm_chain_cd"`
	ShopFranchiseCD     string `sql:"shm_franchise_cd"`
	ProductCd           string `sql:"shm_product_cd"`
	SharedBookStoreCode string `sql:"shm_shared_book_store_code"`
	TelNo               string `sql:"shm_tel_no"`
	BizStartTime        string `sql:"shm_biz_start_time"`
	BizEndTime          string `sql:"shm_biz_end_time"`
	Address             string `sql:"shm_address"`
	ShopNameShort       string `sql:"shm_shop_name_2"`
	IsSelected          bool
}
