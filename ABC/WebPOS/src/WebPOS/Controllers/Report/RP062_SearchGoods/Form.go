package RP062_SearchGoods

type QueryForm struct {
	Page  string `form:"page"`
	Limit int    `form:"limit"`

	GoodsType     []string `form:"goods_type"`
	KeyWord       string   `form:"key_word"`
	KeyWordArrays []string `form:"key_word_arrays"`
	Sort          string   `form:"sort"`
	JanArrays     []string `form:"jan_cd_array"`
	KeySearch     string   `form:"key_search"`

	//Link JAN single good
	JAN       string `form:"jan_cd"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	GroupType string `form:"group_type"`
}
type UserJan struct {
	JanCode        string
	PriorityNumber string
	ProductName    string
	AuthorName     string
	MakerName      string
	UnitPrice      string
	ReleaseDate    string
	Memo           string
}
