package RP060_FavoriteManagement

type QueryForm struct {
	KeySearch          string `form:"key_search"`
	JanCodeList        string `form:"jan_code_list"`
	ProductNameList    string `form:"product_name_list"`
	AuthorNameList     string `form:"author_name_list"`
	MakerNameList      string `form:"publisher_name_list"`
	UnitPriceList      string `form:"usual_price_list"`
	ReleaseDateList    string `form:"release_date_list"`
	LengthListSelected int    `form:"length_list_selected"`
}
