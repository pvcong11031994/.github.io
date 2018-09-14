package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type BQWorkTypeMasterModel struct {
	DB *sql.DB
}

func (this *BQWorkTypeMasterModel) ListWorkType(listWorkType *[]ModelItems.WorkType) error {

	sql := `
SELECT
    mwt_work_type,
    mwt_type_name
FROM
    master_work_type
ORDER BY
    mwt_work_type
    `
	rows, err := this.DB.Query(sql)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	defer rows.Close()

	for rows.Next() {
		w := ModelItems.WorkType{}
		err = rows.Scan(&w.WorkTypeCD, &w.WorkTypeName)
		if err != nil {
			return exterror.WrapExtError(err)
		}
		*listWorkType = append(*listWorkType, w)
	}
	return nil
}
