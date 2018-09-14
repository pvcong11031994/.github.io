package RP059_InitSalesCompare

type QueryForm struct {
	ShopCd           []string `form:"shop_cd"`
	MagazineCdSingle string   `form:"magazine_cd_single"`
	JanArrays        []string `form:"jan_cd_array"`
	JanKey           string   `form:"jan_key"`
	SearchHandleType string   `form:"search_handle_type"`
	ControlType      string   `form:"control_type"`
	SearchDateType   int      `form:"group_type_date"`
	JAN              string   `form:"jan_cd"`

	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	GroupType string `form:"group_type"`

	KeySearch string `form:"key_search"`
}

var ListRank = []string{
	"A",
	"B",
	"C",
	"D",
	"E",
	"F",
	"G",
	"H",
	"I",
	"J",
	"K",
	"L",
	"M",
	"N",
	"O",
	"X",
}
