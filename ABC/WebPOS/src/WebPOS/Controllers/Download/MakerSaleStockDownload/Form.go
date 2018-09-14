package MakerSaleStockDownload

import "encoding/gob"

type Form struct {
	MakerCd  string   `form:"maker_cd"`
	ShopCd   []string `form:"shop_cd"`
	DateFrom string   `form:"date_from"`
	DateTo   string   `form:"date_to"`
	JAN      string   `form:"jan_cd"`
	DataMode string   `form:"data_mode"`
	GoodsType int `form:"bqgm_goods_type"`
}

func init() {

	gob.Register(Form{})
}
