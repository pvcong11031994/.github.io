package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type CorpMasterModel struct {
	DB *sql.DB
}

func (this *CorpMasterModel) GetListCorp(corpCd string, corpItem *[]ModelItems.CorpItem) error {
	sqlString := `
SELECT
	cm_corp_cd,
	cm_corp_name
FROM
	corp_master
`
	rows, err := this.DB.Query(sqlString)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCorpItem := ModelItems.CorpItem{}
		err = rows.Scan(
			&newCorpItem.CorpCd,
			&newCorpItem.CorpName,
		)
		if err != nil {
			return exterror.WrapExtError(err)
		}
		*corpItem = append(*corpItem, newCorpItem)
	}

	return nil
}

func (this *CorpMasterModel) GetListCorpByCorpCd(corpCd string, corpItem *[]ModelItems.CorpItem) error {
	sqlString := `
SELECT
	cm_corp_cd,
	cm_corp_name
FROM
	corp_master
WHERE cm_corp_cd = ?
`
	rows, err := this.DB.Query(sqlString)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCorpItem := ModelItems.CorpItem{}
		err = rows.Scan(
			&newCorpItem.CorpCd,
			&newCorpItem.CorpName,
		)
		if err != nil {
			return exterror.WrapExtError(err)
		}
		*corpItem = append(*corpItem, newCorpItem)
	}

	return nil
}
