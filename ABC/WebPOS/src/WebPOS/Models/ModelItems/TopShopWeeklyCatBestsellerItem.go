package ModelItems

type TopShopWeeklyCatBestsellerItem struct {
	CreateDate  string `sql:"tswcb_create_date"`
	UpdateDate  string `sql:"tswcb_update_date"`
	Servername  string `sql:"tswcb_servername"`
	Dbname      string `sql:"tswcb_dbname"`
	FranchiseCd string `sql:"tswcb_franchise_cd"`
	ChainCd     string `sql:"tswcb_chain_cd"`
	ShopCd      string `sql:"tswcb_shop_cd"`
	Date        string `sql:"tswcb_date"`
	Jan         string `sql:"tswcb_jan"`
	GoodsName   string `sql:"tswcb_goods_name"`
	ArtistName  string `sql:"tswcb_artist_name"`
	MakerName   string `sql:"tswcb_maker_name"`
	Price       string `sql:"tswcb_price"`
	GenreCd     string `sql:"tswcb_genre_cd"`
	GenreName   string `sql:"tswcb_genre_name"`
	SalesCount  int    `sql:"tswcb_sales_count"`
}
