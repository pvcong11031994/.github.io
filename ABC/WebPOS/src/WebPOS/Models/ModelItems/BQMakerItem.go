package ModelItems

type MakerItem struct {
	MakerId    string `sql:"bqmk_makerid"`
	ServerName string `sql:"bqmk_servername"`
	DbName     string `sql:"bqmk_dbname"`
	CreateDate string `sql:"bqmk_create_date"`
	UpdateDate string `sql:"bqmk_update_date"`
	MakerType  string `sql:"bqmk_maker_type"`
	MakerCd    string `sql:"bqmk_maker_cd"`
	MakerName  string `sql:"bqmk_maker_name"`
	LabelCd    string `sql:"bqmk_label_cd"`
	LabelName  string `sql:"bqmk_label_name"`
}
