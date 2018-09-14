package ModelItems

type TopMonthByYearSalesItem struct {
	CreateDate  string `sql:"tmys_create_date"`
	UpdateDate  string `sql:"tmys_update_date"`
	Servername  string `sql:"tmys_servername"`
	Dbname      string `sql:"tmys_dbname"`
	FranchiseCd string `sql:"tmys_franchise_cd"`
	ChainCd     string `sql:"tmys_chain_cd"`
	ShopCd      string `sql:"tmys_shop_cd"`
	Date        string `sql:"tmys_date"`
	Year        string `sql:"tmys_year"`
	Month       string `sql:"tmys_month"`
	Day         string `sql:"tmys_day"`
	GenreCd     string `sql:"tmys_genre_cd"`
	GenreName   string `sql:"tmys_genre_name"`
	SalesCount  int    `sql:"tmys_sales_count"`
	SalesAmount int    `sql:"tmys_sales_amount"`
}

type CompareLastYearByPercentItem struct {
	ShopCd    string `sql:"tmys_shop_cd"`
	ShopName  string `sql:"tmys_shop_name"`
	GenreCd   string `sql:"tmys_genre_cd"`
	GenreName string `sql:"tmys_genre_name"`

	LastDateSalesCount int `sql:"tmys_last_date_sales_count"`
	LastYearSalesCount int `sql:"tmys_last_year_sales_count"`
	PercentSalesCount  string

	LastDateSalesAmount int `sql:"tmys_last_date_sales_amount"`
	LastYearSalesAmount int `sql:"tmys_last_year_sales_amount"`
	PercentSalesAmount  string
}
