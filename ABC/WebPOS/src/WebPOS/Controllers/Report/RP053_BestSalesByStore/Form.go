package RP053_BestSalesByStore

type QueryForm struct {
	LayoutColArr []string
	LayoutRowArr []string
	LayoutSumArr []string
	LayoutCols   string `form:"layout_cols"`
	LayoutRows   string `form:"layout_rows"`
	LayoutSums   string `form:"layout_sums"`
	ShopCd       string `form:"shop_cd"`

	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	MonthFrom string `form:"month_from"`
	MonthTo   string `form:"month_to"`
	WeekFrom  string `form:"week_from"`
	WeekTo    string `form:"week_to"`
	GroupType string `form:"group_type"`

	MediaGroup1Cd []string `form:"media_group1_cd"`
	MediaGroup2Cd []string `form:"media_group2_cd"`
	MediaGroup3Cd []string `form:"media_group3_cd"`

	MakerCd          []string `form:"maker_cd"`
	Page             int      `form:"page"`
	SearchHandleType string   `form:"search_handle_type"`
}
