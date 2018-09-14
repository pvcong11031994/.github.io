package RP052_ShopTotalSum

type QueryForm struct {
	ShopCd []string `form:"shop_cd"`

	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	MonthFrom string `form:"month_from"`
	MonthTo   string `form:"month_to"`
	WeekFrom  string `form:"week_from"`
	WeekTo    string `form:"week_to"`
	GroupType string `form:"group_type"`
	DateRank  int    `form:"date_rank"`

	MediaGroup1Cd []string `form:"media_group1_cd"`
	MediaGroup2Cd []string `form:"media_group2_cd"`
	MediaGroup3Cd []string `form:"media_group3_cd"`

	MakerCd          []string `form:"maker_cd"`
	JAN              string   `form:"jan_cd"`
	Page             int64    `form:"page"`
	Limit            int64    `form:"limit"`
	SearchHandleType string   `form:"search_handle_type"`
}
