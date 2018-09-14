package Models

import (
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type BQPosMasterModel struct {
	DB *sql.DB
}

func (this *BQPosMasterModel) GetPos(shop string) ([]string, error) {
	var list_pos []string
	strQuery := `
SELECT bqpm_pos_no
FROM bq_pos_master
WHERE CONCAT(bqpm_servername,'|',bqpm_shop_cd) = ?
	`
	rows, err := this.DB.Query(strQuery, shop)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		pos_no := ""
		err = rows.Scan(&pos_no)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		list_pos = append(list_pos, pos_no)
	}
	return list_pos, nil
}
