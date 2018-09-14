package ModelItems

type TopDayGenreSalesItem struct {
	TrnDate     string `sql:"tpgs_date"`
	TrnYear     string `sql:"tpgs_year"`
	TrnMonth    string `sql:"tpgs_month"`
	TrnDay      string `sql:"tpgs_day"`
	SalesAmount int64  `sql:"tpgs_sales_amount"`
	GenreCd     string `sql:"tpgs_genre_cd"`
	GenreName   string `sql:"tpgs_genre_name"`
}
