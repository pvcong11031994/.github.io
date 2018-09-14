package ModelItems

type TopMonthlySalesItem struct {
	TrnDate     string `sql:"tms_date"`
	TrnYear     string `sql:"tms_year"`
	TrnMonth    string `sql:"tms_month"`
	TrnDay      string `sql:"tms_day"`
	SalesAmount int64  `sql:"tms_sales_amount"`
}
