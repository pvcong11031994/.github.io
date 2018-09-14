package RP061_ShopSales

type QueryForm struct {
	ShopCd    []string `form:"shop_cd"`
	DateFrom  string   `form:"date_from"`
	DateTo    string   `form:"date_to"`
	MonthFrom string   `form:"month_from"`
	MonthTo   string   `form:"month_to"`
	WeekFrom  string   `form:"week_from"`
	WeekTo    string   `form:"week_to"`
	GroupType string   `form:"group_type"`
	Page      int      `form:"page"`
	Limit     int      `form:"limit"`

	JanArrays        []string `form:"jan_cd_array"`
	JanSingle        string   `form:"jan_cd_single"`
	SearchHandleType string   `form:"search_handle_type"`
	MakerCd          []string `form:"maker_code"`
	LinkRevert       string   `form:"link_revert"`
	JAN              string   `form:"jan_cd"`
	ControlType      string   `form:"control_type"`
	DownloadType     string   `form:"download_type"`
}
