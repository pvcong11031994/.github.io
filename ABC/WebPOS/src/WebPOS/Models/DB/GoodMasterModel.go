package Models

import (
	"WebPOS/Common"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
)

type GoodMasterModel struct {
	DB *sql.DB
}

type GoodMasterItem struct {
	Jan             string `sql:"gm_jan"`
	GoodsName       string `sql:"gm_goods_name"`
	Author          string `sql:"gm_author"`
	PublisherCd     string `sql:"gm_publisher_cd"`
	PublisherName   string `sql:"gm_publisher_name"`
	Price           string `sql:"gm_price"`
	PublishDate     string `sql:"gm_publish_date"`
	BGenreCd        string `sql:"gm_media_cd"`
	MediaName       string `sql:"gm_media_name"`
	FirstInputCount string `sql:"gm_first_input_count"`
	LocationCd      string `sql:"location_cd"`
}

func (this *GoodMasterModel) ListByJans(jan_cd, shop []string) ([]GoodMasterItem, error) {
	listResult := []GoodMasterItem{}
	query := `
SELECT
	gm_jan,
	gm_goods_name,
	gm_author,
	gm_publisher_cd,
	gm_publisher_name,
	gm_price,
	gm_published_date,
	gm_media_cd,
	gm_media_name,
	gm_first_input_count,
	location_cd
FROM
    ao_goods_master
LEFT JOIN
(
	SELECT
		bqls_jan_cd jan,
		bqls_tana_cd location_cd
	FROM
		bq_location_setting
	WHERE
	    bqls_tana_type = '1'
	AND bqls_jan_cd IN (` + Common.SQLPara(jan_cd) + `)
	AND CONCAT(bqls_server_name, bqls_shop_cd) IN (` + Common.SQLPara(shop) + `)
) bqls
 ON gm_jan = jan
WHERE
    gm_jan IN (` + Common.SQLPara(jan_cd) + `)`

	args := []interface{}{}
	args = append(args, Common.ToInterfaceArray(jan_cd)...)
	args = append(args, Common.ToInterfaceArray(shop)...)
	args = append(args, Common.ToInterfaceArray(jan_cd)...)

	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newItem := GoodMasterItem{}
		err = db.SqlScanStruct(rows, &newItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listResult = append(listResult, newItem)
	}
	return listResult, nil
}
