package ModelItems

type UserJanCodeItem struct {
	CreateDate     string `sql:"ujc_create_datetime"`
	UpdateDate     string `sql:"ujc_update_datetime"`
	UserId         string `sql:"ujc_user_id"`
	JanCode        string `sql:"ujc_jan_code"`
	ProductName    string `sql:"ujc_product_name"`
	MakerName      string `sql:"ujc_maker_name"`
	AuthorName     string `sql:"ujc_author_name"`
	SellingDate    string `sql:"ujc_selling_date"`
	ListPrice      string `sql:"ujc_list_price"`
	Memo           string `sql:"ujc_user_jan_inf"`
	PriorityNumber string `sql:"ujc_priority_number"`
}
