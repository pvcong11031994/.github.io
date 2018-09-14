package ModelItems

type TopShopBestsellerItem struct {
	JanCd       string `sql:"tsb_jan"`
	GenreCd     string `sql:"tsb_genre_cd"`
	GenreName   string `sql:"tsb_genre_name"`
	GoodsName   string `sql:"tsb_goods_name"`
	ArtistName  string `sql:"tsb_artist_name"`
	MakerName   string `sql:"tsb_maker_name"`
	Price       int64  `sql:"tsb_price"`
	SalesCount  int64  `sql:"tsb_sales_count"`
	SalesAmount int64  `sql:"tsb_sales_amount"`
	Stock       int64  `sql:"tsb_stock_count"`
}
