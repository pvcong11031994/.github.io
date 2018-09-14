package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
)

type BQLocationMasterModel struct {
	DB *sql.DB
}

func (this *BQLocationMasterModel) GetListLocationByUser(userId string) ([]ModelItems.BQLocationItem, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listLocation := []ModelItems.BQLocationItem{}

	sqlVar := []interface{}{}

	sqlString := `
SELECT
   bqlm.bqlm_tana_cd,
   bqlm.bqlm_tana_name
FROM
   bq_location_setting
LEFT JOIN bq_location_master bqlm
ON bqlm_uid = bqls_location_id
AND bqlm_shop_cd = bqls_shop_cd
WHERE bqls_tana_type = '1'
`
	if !role.IsHonbu || role.ShopCD != "" {
		sqlString += ` AND bqlm_shop_cd = ? `
		sqlVar = append(sqlVar, role.ShopCD)
	} else {
		sqlString += ` AND FALSE `
	}

	sqlString += `
GROUP BY
    bqlm_tana_cd,
   bqlm_tana_name
ORDER BY
    bqlm_tana_cd

   `
	rows, err := this.DB.Query(sqlString, sqlVar...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newBQLocationItem := ModelItems.BQLocationItem{}
		err = db.SqlScanStruct(rows, &newBQLocationItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listLocation = append(listLocation, newBQLocationItem)
	}

	return listLocation, nil
}
