package ModelItems

type SettingShopReferItem struct {
	CreateDate       string `sql:"ssr_create_date"`
	UpdateDate       string `sql:"ssr_update_date"`
	ShopCd           string `sql:"ssr_shop_cd"`
	ReferGenreMaster string `sql:"ssr_refer_genre_master"`
	ReferGenreType   string `sql:"ssr_refer_genre_type"`
	StartGenreCd     string `sql:"ssr_start_genre_cd"`
	EndGenreCd       string `sql:"ssr_end_genre_cd"`
	GenreCdCheckOff  string `sql:"ssr_genre_cd_check_off"`
	ToritsugiCd1     string `sql:"ssr_toritsugi_cd1"`
	ToritsugiCd2     string `sql:"ssr_toritsugi_cd2"`
	ToritsugiCd3     string `sql:"ssr_toritsugi_cd3"`
	PriorityGm       string `sql:"ssr_priority_gm"`
}
