package ModelItems

type FranchiseCdItem struct {
	FranchiseCd   string `sql:"shm_franchise_cd"`
	FranchiseName string `sql:"shm_franchise_name"`
}

type FranchiseItem struct {
	FranchiseCd   string `sql:"sfgm_franchise_cd"`
	FranchiseName string `sql:"sfgm_franchise_name"`
	FranchiseShop []*FranchiseShopItem
}

type FranchiseShopItem struct {
	ServerName        string `sql:"shm_server_name"`
	ShopCd            string `sql:"shm_shop_cd"`
	ShopName          string `sql:"shm_shop_name"`
	FranchiseShopUser []*FranchiseShopUserItem
}

type FranchiseShopUserItem struct {
	UserID   string
	UserName string
}
