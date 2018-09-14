package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
	"strings"
)

type BQChargingSettingModel struct {
	DB *sql.DB
}

func (this *BQChargingSettingModel) CheckExistCharging(serverName, shopCd, userCd string) (bool, error) {
	query := `
SELECT
	IF(COUNT(*), 'true', 'false')
FROM
	bq_charging_setting
WHERE
	bqcs_server_name = ?
	AND bqcs_shop_cd = ?
 	AND bqcs_user_ID = ?
`

	args := []interface{}{}
	args = append(args, serverName)
	args = append(args, shopCd)
	args = append(args, userCd)

	var isExist bool
	err := this.DB.QueryRow(query, args...).Scan(&isExist)
	if err != nil {
		return false, exterror.WrapExtError(err)
	}

	return isExist, nil
}

func (this *BQChargingSettingModel) RegisterCharging(list []ModelItems.BQChargingSettingItem) error {

	strValue := strings.Repeat("(?, ?, ?, ?, ?),", len(list))
	args := []interface{}{}
	for _, item := range list {
		args = append(args, item.ServerName)
		args = append(args, item.ShopCd)
		args = append(args, item.UserID)
		args = append(args, item.Bairitsu)
		args = append(args, item.DeleteFlag)
	}
	strValue = strings.TrimRight(strValue, ",")
	query := (`
INSERT INTO bq_charging_setting (
	 bqcs_server_name
	,bqcs_shop_cd
	,bqcs_user_ID
	,bqcs_bairitsu
	,bqcs_delete_flag)
VALUES
` + strValue + `
ON DUPLICATE KEY UPDATE
	bqcs_bairitsu	 		=VALUES(bqcs_bairitsu)
	,bqcs_delete_flag	 		=VALUES(bqcs_delete_flag)
	`)
	_, err := this.DB.Exec(query, args...)

	return exterror.WrapExtError(err)
}

func (this *BQChargingSettingModel) SelectCharging(serverShop, franchiseCd []string, franchiseGroupCd string) ([]ModelItems.BQChargingSettingItem, error) {
	listCharging := []ModelItems.BQChargingSettingItem{}
	strWhere := ""
	strValue := ""
	args := []interface{}{}
	if len(serverShop) > 0 {
		strValue = strings.Repeat("?,", len(serverShop))
		for _, item := range serverShop {
			args = append(args, item)
		}
		strWhere = " CONCAT(bqcs_server_name ,'-',bqcs_shop_cd)  IN "
	} else {
		strValue = strings.Repeat("?,", len(franchiseCd))
		for _, item := range franchiseCd {
			args = append(args, item)
		}
		strWhere = " shm_franchise_cd  IN "
	}

	strValue = "(" + strings.TrimRight(strValue, ",") + ")"
	args = append(args, franchiseGroupCd)
	query := `
SELECT
	bqcs_server_name
	,bqcs_shop_cd
	,shm_shop_name
	,shm_franchise_cd
	,shm_franchise_name
	,IFNULL(bqcs_user_ID,'')
	,IFNULL(um_user_name,'')
	,bqcs_bairitsu
	,bqcs_delete_flag
FROM bq_charging_setting
LEFT JOIN shop_master_show
	ON bqcs_server_name = shm_server_name
	AND bqcs_shop_cd = shm_shop_cd
LEFT JOIN ao_user_master
	ON bqcs_user_ID = um_user_ID
LEFT JOIN shop_franchise_group_master
	ON shm_franchise_cd = sfgm_franchise_cd
WHERE
	 ` + strWhere + strValue + `
	 AND sfgm_franchise_group_cd = ?
	`
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return listCharging, exterror.WrapExtError(err)
	}
	defer rows.Close()

	for rows.Next() {
		g := ModelItems.BQChargingSettingItem{}
		err = rows.Scan(&g.ServerName,
			&g.ShopCd,
			&g.ShopName,
			&g.FranchiseCd,
			&g.FranchiseName,
			&g.UserID,
			&g.UserName,
			&g.Bairitsu,
			&g.DeleteFlag)
		if err != nil {
			return listCharging, exterror.WrapExtError(err)
		}
		listCharging = append(listCharging, g)
	}
	return listCharging, nil
}

func (this *BQChargingSettingModel) GetInfoChargingByCd(shopCd, serverName, userId string) (ModelItems.BQChargingSettingItem, error) {

	infoCharging := ModelItems.BQChargingSettingItem{}
	args := []interface{}{}
	args = append(args, serverName)
	args = append(args, shopCd)
	args = append(args, userId)
	query := `
SELECT
	bqcs_server_name,
	bqcs_shop_cd,
	bqcs_user_ID,
	IFNULL(bqcs_bairitsu,0)
FROM bq_charging_setting
WHERE
	bqcs_delete_flag = '0'
	AND bqcs_server_name = ?
	AND bqcs_shop_cd = ?
	AND bqcs_user_ID = ?
`
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return infoCharging, exterror.WrapExtError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&infoCharging.ServerName,
			&infoCharging.ShopCd,
			&infoCharging.UserID,
			&infoCharging.Bairitsu,
		)
		if err != nil {
			return infoCharging, exterror.WrapExtError(err)
		}
	}
	return infoCharging, nil
}
