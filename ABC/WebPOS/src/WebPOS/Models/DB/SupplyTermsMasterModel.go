package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type SupplyTermsMasterModel struct {
	DB *sql.DB
}

func (this *SupplyTermsMasterModel) ListSupplyTerms(listSupplyTerms *[]ModelItems.SupplyTerms) error {

	sql := `
SELECT
    stm_supply_terms_cd,
    stm_supply_terms_name
FROM
    supply_terms_master
ORDER BY
    stm_supply_terms_cd
    `
	rows, err := this.DB.Query(sql)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	defer rows.Close()

	for rows.Next() {
		w := ModelItems.SupplyTerms{}
		err = rows.Scan(&w.SupplyTermsCD, &w.SupplyTermsName)
		if err != nil {
			return exterror.WrapExtError(err)
		}
		*listSupplyTerms = append(*listSupplyTerms, w)
	}
	return nil
}
