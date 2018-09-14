package RP058_SalesComparison

type QueryForm struct {
	ShopCd    []string `form:"shop_cd"`
	DateFrom  string   `form:"date_from"`
	DateTo    string   `form:"date_to"`
	MonthFrom string   `form:"month_from"`
	MonthTo   string   `form:"month_to"`
	WeekFrom  string   `form:"week_from"`
	WeekTo    string   `form:"week_to"`
	GroupType string   `form:"group_type"`

	JanArrays        []string `form:"jan_cd_array"`
	JanKey           string   `form:"jan_key"`
	SearchHandleType string   `form:"search_handle_type"`
	MakerCd          []string `form:"maker_code"`
	LinkRevert       string   `form:"link_revert"`
	FlagJan          string   `form:"flag_jan"`
	FlagSingleItem   string   `form:"flag_single_item"`
	JAN              string   `form:"jan_cd"`
	KeySearch        string   `form:"key_search"`
}
