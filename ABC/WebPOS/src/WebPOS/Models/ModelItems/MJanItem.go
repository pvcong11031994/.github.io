package ModelItems

type MJanItem struct {
	JanCode     string `sql:"jan_code"`
	ProductName string `sql:"product_name"`
	AuthorName  string `sql:"author_name"`
	MakerName   string `sql:"maker_name"`
	ListPrice   int64  `sql:"list_price"`
	SellingDate string `sql:"selling_date"`
	StockInf    string `sql:"stock_inf"`
}
