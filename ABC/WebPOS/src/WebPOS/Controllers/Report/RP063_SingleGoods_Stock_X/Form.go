package RP063_SingleGoods_Stock_X

type QueryFormSingleGoods struct {
	ShopCd    []string `form:"shop_cd"`
	DateFrom  string   `form:"date_from"`
	DateTo    string   `form:"date_to"`
	MonthFrom string   `form:"month_from"`
	MonthTo   string   `form:"month_to"`
	WeekFrom  string   `form:"week_from"`
	WeekTo    string   `form:"week_to"`
	GroupType string   `form:"group_type"`

	JAN              string `form:"jan_code"`
	FlagSingleItem   string `form:"flag_single_item"`
	SearchType       string `form:"search_type"`
	Page             int    `form:"page"`
	Limit            int    `form:"limit"`
	SearchHandleType string `form:"search_handle_type"`

	MediaGroup1Cd []string `form:"media_group1_cd"`
	MediaGroup2Cd []string `form:"media_group2_cd"`
	MediaGroup3Cd []string `form:"media_group3_cd"`
	MediaGroup4Cd []string `form:"media_group4_cd"`

	JanMakerCode        []string `form:"jan_maker_code"`
	MakerCd             []string `form:"maker_code"`
	MagazineCd          []string `form:"magazine_cd"`
	MagazineCodeMonth   string   `form:"magazine_code_month"`
	MagazineCodeWeek    string   `form:"magazine_code_week"`
	MagazineCodeQuarter string   `form:"magazine_code_quarter"`
	ControlType         string   `form:"control_type"`
	DownloadType        string   `form:"download_type"`
	LinkRevert          string   `form:"link_revert"`
	// Shop Sales
	JanArrays      []string `form:"jan_cd_array"`
	JanSingle      string   `form:"jan_cd_single"`
	SearchDateType int      `form:"group_type_date"`

	//Init_shop_sales
	MagazineCdSingle string `form:"magazine_cd_single"`
}
