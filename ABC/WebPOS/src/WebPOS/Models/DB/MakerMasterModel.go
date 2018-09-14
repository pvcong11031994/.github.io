package Models

import (
	"database/sql"
	"WebPOS/Models/ModelItems"
	"github.com/goframework/gf/exterror"
)

type MakerMasterModel struct {
	DB *sql.DB
}

func (this *MakerMasterModel) GetPublisherInfoByUser(userId string) (maker ModelItems.MakerItem, err error) {

	sqlString := `
SELECT
	bqmk_makerid,
	bqmk_servername,
	bqmk_dbname,
	bqmk_create_date,
	bqmk_update_date,
	bqmk_maker_type,
	bqmk_maker_cd,
	bqmk_maker_name,
	bqmk_label_cd,
	bqmk_label_name
FROM
	ao_user_master um
LEFT JOIN bq_maker bqmk
ON
	um.um_server_name = bqmk.bqmk_servername
	AND um.um_dept_cd = bqmk.bqmk_maker_cd
WHERE
	um.um_user_ID = ?
`

	err = this.DB.QueryRow(sqlString, userId).Scan(
		&maker.MakerId,
		&maker.ServerName,
		&maker.DbName,
		&maker.CreateDate,
		&maker.UpdateDate,
		&maker.MakerType,
		&maker.MakerCd,
		&maker.MakerName,
		&maker.LabelCd,
		&maker.LabelName,
	)
	if err != nil {
		return ModelItems.MakerItem{}, exterror.WrapExtError(err)
	}

	return maker, nil
}