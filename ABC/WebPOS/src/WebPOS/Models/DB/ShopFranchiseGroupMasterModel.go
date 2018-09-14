package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
)

type ShopFranchiseGroupMasterModel struct {
	DB *sql.DB
}

func (this *ShopFranchiseGroupMasterModel) GetFranchiseGroups(listGroup []string) ([]ModelItems.FranchiseGroupItem, error) {

	sql := `
SELECT
	sfgm_franchise_group_cd,
	sfgm_franchise_group_name
FROM shop_franchise_group_master
WHERE
	sfgm_franchise_group_cd IN (` + Common.SQLPara(listGroup) + `)
GROUP BY
	sfgm_franchise_group_cd
`

	rows, err := this.DB.Query(sql, Common.ToInterfaceArray(listGroup)...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	list := []ModelItems.FranchiseGroupItem{}
	for rows.Next() {
		item := ModelItems.FranchiseGroupItem{}
		err = db.SqlScanStruct(rows, &item)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		list = append(list, item)
	}

	return list, nil
}

func (this *ShopFranchiseGroupMasterModel) GetListFranchiseCD(franchiseGroupCD string, listFranchiseCd *[]ModelItems.FranchiseCdItem) error {
	sqlString := `
SELECT
	sfgm_franchise_cd
FROM
	shop_franchise_group_master
WHERE
	sfgm_franchise_group_cd = ?
`

	rows, err := this.DB.Query(sqlString, franchiseGroupCD)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		g := ModelItems.FranchiseCdItem{}
		err = rows.Scan(&g.FranchiseCd)
		if err != nil {
			return exterror.WrapExtError(err)
		}
		*listFranchiseCd = append(*listFranchiseCd, g)
	}

	return nil
}

func (this *ShopFranchiseGroupMasterModel) GetNameFranchise(franchiseGroupCd, franchiseCd string) (string, string, error) {
	franchiseGroupName := ""
	franchiseName := ""
	sqlString := `
SELECT
	sfgm_franchise_group_name,
	sfgm_franchise_name
FROM
	shop_franchise_group_master
WHERE
	sfgm_franchise_group_cd = ?
	AND sfgm_franchise_cd = ?
`

	rows, err := this.DB.Query(sqlString, franchiseGroupCd, franchiseCd)
	if err != nil {
		return "", "", exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&franchiseGroupName, &franchiseName)
		if err != nil {
			return "", "", exterror.WrapExtError(err)
		}
	}

	return franchiseGroupName, franchiseName, nil
}
