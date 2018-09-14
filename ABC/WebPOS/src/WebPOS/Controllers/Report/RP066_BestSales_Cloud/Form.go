package RP066_BestSales_Cloud

type QueryForm struct {
	LayoutColArr []string
	LayoutRowArr []string
	LayoutSumArr []string
	LayoutCols   string   `form:"layout_cols"`
	LayoutRows   string   `form:"layout_rows"`
	LayoutSums   string   `form:"layout_sums"`
	ShopCd       []string `form:"shop_cd"`
	ShopName     []string `form:"shop_name"`

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
	MediaGroup4Cd []string `form:"media_group4_cd"`

	JanMakerCode        []string `form:"jan_maker_code"`
	MakerCd             []string `form:"maker_code"`
	MagazineCd          []string `form:"magazine_cd"`
	JAN                 string   `form:"jan_cd"`
	FlagSingleItem      string   `form:"flag_single_item"`
	SearchType          string   `form:"search_type"`
	Page                int      `form:"page"`
	MagazineCodeMonth   string   `form:"magazine_code_month"`
	MagazineCodeWeek    string   `form:"magazine_code_week"`
	MagazineCodeQuarter string   `form:"magazine_code_quarter"`
	Limit               int      `form:"limit"`
	SearchHandleType    string   `form:"search_handle_type"`
	ControlType         string   `form:"control_type"`
	DownloadType        string   `form:"download_type"`
}
