package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"strings"
)

type BQLogSearchChargingModel struct {
	DB *sql.DB
}

// get list data group by shop_cd, server_name
func (this *BQLogSearchChargingModel) GetListBySearchTime(condition map[string]interface{}) ([]ModelItems.BQLogSearchChargingItem, int64, int64, error) {

	targetMonth := condition["target_month"].(string)
	targetDateForm := condition["target_date_form"].(string)
	targetDateTo := condition["target_date_to"].(string)
	listFranchiseCd := condition["franchise_cd"].([]string)
	listShopCd := condition["shop_cd"].([]string)

	conditionWhere := ""
	args := []interface{}{}
	args = append(args, condition["franchise_group_cd"].(string))

	if targetDateForm != "" || targetDateTo != "" {
		if targetDateForm != "" {
			conditionWhere += `
			AND LEFT(REPLACE(bqls_searching_date, '-', ''), 8) >= ?
			`
			args = append(args, targetDateForm)
		}
		if targetDateTo != "" {
			conditionWhere += `
			AND LEFT(REPLACE(bqls_searching_date, '-', ''), 8) <= ?
			`
			args = append(args, targetDateTo)
		}
	} else if targetMonth != "" {
		conditionWhere += `
		AND LEFT(REPLACE(bqls_searching_date, '-', ''), 6) = ?
		`
		args = append(args, targetMonth)
	}

	if len(listFranchiseCd) > 0 {
		conditionWhere += `
		AND bqls_franchise_cd IN (` + strings.Trim(strings.Repeat("?,", len(listFranchiseCd)), ",") + `)
		`
		for _, value := range listFranchiseCd {
			args = append(args, value)
		}
	}

	if len(listShopCd) > 0 {
		conditionWhere += `
		AND CONCAT(bqls_server_name,'-',bqls_shop_cd) IN (` + strings.Trim(strings.Repeat("?,", len(listShopCd)), ",") + `)
		`
		for _, value := range listShopCd {
			args = append(args, value)
		}
	}

	sqlStrings := `
SELECT
	SUM(bqls_GCP_charging) AS sum_GCP_charging,
	SUM(bqls_VJ_charging) AS sum_VJ_charging,
	bqls_shop_cd,
	shm_shop_name,
	bqls_server_name,
	bqls_franchise_cd,
	sfgm_franchise_name,
	COUNT(bqls_searching_date) count_searching_date,
	IFNULL(bqcs.bairitsu,'') AS bqcs_bairitsu
FROM bq_log_search_charging
LEFT JOIN (
	SELECT
		bqcs_bairitsu AS bairitsu,
		bqcs_server_name,
		bqcs_shop_cd
	FROM
		bq_charging_setting
	GROUP BY
		bqcs_server_name,
		bqcs_shop_cd) bqcs
	ON bqcs.bqcs_server_name = bqls_server_name
	AND bqcs.bqcs_shop_cd = bqls_shop_cd
LEFT JOIN shop_franchise_group_master
	ON sfgm_franchise_cd = bqls_franchise_cd
LEFT JOIN shop_master_show
	ON shm_server_name = bqls_server_name
	AND shm_shop_cd = bqls_shop_cd
WHERE
	sfgm_franchise_group_cd = ?
	` + conditionWhere + `
GROUP BY
	bqls_shop_cd,
	shm_shop_name,
	bqls_server_name,
	bqls_franchise_cd,
	sfgm_franchise_name,
	bqcs_bairitsu
`

	rows, err := this.DB.Query(sqlStrings, args...)
	if err != nil {
		return nil, 0, 0, exterror.WrapExtError(err)
	}
	defer rows.Close()

	resultList := []ModelItems.BQLogSearchChargingItem{}
	var sumGCPAllShop int64 = 0
	var sumVJAllShop int64 = 0
	for rows.Next() {
		item := ModelItems.BQLogSearchChargingItem{}
		err = rows.Scan(
			&item.SumGCPCharging,
			&item.SumVJCharging,
			&item.ShopCd,
			&item.ShopName,
			&item.ServerName,
			&item.FranchiseCd,
			&item.FranchiseName,
			&item.CountSearchingDate,
			&item.Bairitsu)
		if err != nil {
			return nil, 0, 0, exterror.WrapExtError(err)
		}
		sumGCPAllShop = sumGCPAllShop + item.SumGCPCharging
		sumVJAllShop = sumVJAllShop + item.SumVJCharging
		resultList = append(resultList, item)
	}

	return resultList, sumGCPAllShop, sumVJAllShop, nil
}

// get list detail group by all key
func (this *BQLogSearchChargingModel) GetListDetail(condition map[string]interface{}) ([]ModelItems.BQLogSearchChargingItem, int64, int64, error) {

	targetMonth := condition["target_month"].(string)
	targetDateForm := condition["target_date_form"].(string)
	targetDateTo := condition["target_date_to"].(string)
	franchiseCd := condition["franchise_cd"].(string)
	shopCd := condition["shop_cd"].(string)
	serverName := condition["server_name"].(string)

	conditionWhere := ""
	args := []interface{}{}
	args = append(args, condition["franchise_group_cd"].(string))

	if targetDateForm != "" || targetDateTo != "" {
		if targetDateForm != "" {
			conditionWhere += `
			AND LEFT(REPLACE(bqls_searching_date, '-', ''), 8) >= ?
			`
			args = append(args, targetDateForm)
		}
		if targetDateTo != "" {
			conditionWhere += `
			AND LEFT(REPLACE(bqls_searching_date, '-', ''), 8) <= ?
			`
			args = append(args, targetDateTo)
		}
	} else if targetMonth != "" {
		conditionWhere += `
		AND LEFT(REPLACE(bqls_searching_date, '-', ''), 6) = ?
		`
		args = append(args, targetMonth)
	}

	if len(franchiseCd) > 0 {
		conditionWhere += `
		AND bqls_franchise_cd = ?
		`
		args = append(args, franchiseCd)
	}

	if len(shopCd) > 0 && len(serverName) > 0 {
		conditionWhere += `
		AND bqls_server_name = ?
		AND bqls_shop_cd = ?
		`
		args = append(args, serverName)
		args = append(args, shopCd)
	}

	sqlStrings := `
SELECT
	SUM(bqls_GCP_charging) AS sum_GCP_charging,
	SUM(bqls_VJ_charging) AS sum_VJ_charging,
	LEFT(REPLACE(bqls_searching_date,'-',''), 8),
	bqls_user_ID,
	bqls_user_name,
	IFNULL(bqls_use_menu, ''),
	bqls_use_TBL_size
FROM bq_log_search_charging
LEFT JOIN shop_franchise_group_master
	ON sfgm_franchise_cd = bqls_franchise_cd
WHERE
	sfgm_franchise_group_cd = ?
	` + conditionWhere + `
GROUP BY
	bqls_searching_date,
	bqls_user_ID,
	bqls_user_name,
	bqls_use_menu,
	bqls_use_TBL_size
ORDER BY
	bqls_searching_date DESC
`
	rows, err := this.DB.Query(sqlStrings, args...)
	if err != nil {
		return nil, 0, 0, exterror.WrapExtError(err)
	}
	defer rows.Close()

	resultList := []ModelItems.BQLogSearchChargingItem{}
	var sumGCPAllShop int64 = 0
	var sumVJAllShop int64 = 0
	for rows.Next() {
		item := ModelItems.BQLogSearchChargingItem{}
		err = rows.Scan(
			&item.SumGCPCharging,
			&item.SumVJCharging,
			&item.SearchingDate,
			&item.UserID,
			&item.UserName,
			&item.UserMenu,
			&item.UseTBLSize,
		)
		if err != nil {
			return nil, 0, 0, exterror.WrapExtError(err)
		}
		sumGCPAllShop = sumGCPAllShop + item.SumGCPCharging
		sumVJAllShop = sumVJAllShop + item.SumVJCharging
		resultList = append(resultList, item)
	}

	return resultList, sumGCPAllShop, sumVJAllShop, nil
}

type ChanChargesSearch struct {
	Data *ModelItems.DataChargesByUserItem
	Err  error
}

func (this *BQLogSearchChargingModel) ListChargesByUser(condition map[string]interface{}) chan ChanChargesSearch {

	args := []interface{}{}
	strCondition := ""
	// 検索日指定
	if condition["search_month"] != nil {
		strCondition += " AND DATE_FORMAT(bqls_searching_date,'%Y%m') = ? "
		args = append(args, condition["search_month"].(string))
	} else {
		if condition["date_from"] != nil {
			strCondition += " AND DATE_FORMAT(bqls_searching_date,'%Y%m%d') >= ? "
			args = append(args, condition["date_from"])
		}
		if condition["date_to"] != nil {
			strCondition += " AND DATE_FORMAT(bqls_searching_date,'%Y%m%d') <= ? "
			args = append(args, condition["date_to"])
		}
	}
	// フランチャイズ
	if condition["franchise_cd"] != nil {
		strCondition += " AND bqls_franchise_cd IN (" + Common.SQLPara(condition["franchise_cd"].([]string)) + ") "
		for _, franchise_cd := range condition["franchise_cd"].([]string) {
			args = append(args, franchise_cd)
		}
	}
	// 店舗
	if condition["shop_cd"] != nil {
		strCondition += " AND CONCAT(bqls_server_name,'|',bqls_shop_cd) IN (" + Common.SQLPara(condition["shop_cd"].([]string)) + ") "
		for _, shop_cd := range condition["shop_cd"].([]string) {
			args = append(args, shop_cd)
		}
	}

	sqlQuery := `
SELECT
	bqls_server_name        AS bqls_server_name,
	bqls_shop_cd            AS bqls_shop_cd_list,
	bqls_shop_name          AS bqls_shop_name_list,
	COUNT(bqls_shop_cd)     AS bqls_total_search_list,
	SUM(bqls_VJ_charging)   AS bqls_total_charges_list
FROM
	bq_log_search_charging
WHERE TRUE
` + strCondition + `
GROUP BY
	bqls_server_name,
	bqls_shop_cd
ORDER BY
	bqls_shop_cd
`

	outputChan := make(chan ChanChargesSearch)
	go func() {
		rows, err := this.DB.Query(sqlQuery, args...)
		if err != nil {
			outputChan <- ChanChargesSearch{nil, exterror.WrapExtError(err)}
			close(outputChan)
			return
		}
		defer rows.Close()
		for rows.Next() {
			item := ModelItems.DataChargesByUserItem{}
			err = db.SqlScanStruct(rows, &item)
			if err != nil {
				outputChan <- ChanChargesSearch{nil, exterror.WrapExtError(err)}
				close(outputChan)
				return
			}
			item.SearchDateList = condition["search_date"].(string)
			outputChan <- ChanChargesSearch{&item, nil}
		}
		close(outputChan)
	}()
	return outputChan
}

func (this *BQLogSearchChargingModel) DetailChargesByUser(condition map[string]interface{}) chan ChanChargesSearch {

	args := []interface{}{}
	strCondition := ""
	// 検索日指定
	if condition["search_month"] != nil {
		strCondition += " AND DATE_FORMAT(bqls_searching_date,'%Y%m') = ? "
		args = append(args, condition["search_month"])
	} else {
		if condition["date_from"] != nil {
			strCondition += " AND DATE_FORMAT(bqls_searching_date,'%Y%m%d') >= ? "
			args = append(args, condition["date_from"])
		}
		if condition["date_to"] != nil {
			strCondition += " AND DATE_FORMAT(bqls_searching_date,'%Y%m%d') <= ? "
			args = append(args, condition["date_to"])
		}
	}
	// 店舗
	strCondition += " AND CONCAT(bqls_server_name,'|',bqls_shop_cd) = ? "
	args = append(args, condition["shop_cd"])

	sqlQuery := `
SELECT
	DATE_FORMAT(bqls_searching_date,'%Y%m%d')       AS bqls_searching_date,
	bqls_shop_cd                                    AS bqls_shop_cd_detail,
	bqls_shop_name                                  AS bqls_shop_name_detail,
	bqls_user_ID                                    AS bqls_user_ID_detail,
	bqls_user_name                                  AS bqls_user_name_detail,
	bqls_use_menu                                   AS bqls_use_menu,
	SUM(bqls_VJ_charging)                           AS bqls_total_charges_detail
FROM
	bq_log_search_charging
WHERE TRUE
` + strCondition + `
GROUP BY
	DATE_FORMAT(bqls_searching_date,'%Y%m%d'),
	bqls_user_ID,
	bqls_use_menu
ORDER BY
	bqls_searching_date,
	bqls_user_ID
`

	outputChan := make(chan ChanChargesSearch)
	go func() {
		rows, err := this.DB.Query(sqlQuery, args...)
		if err != nil {
			outputChan <- ChanChargesSearch{nil, exterror.WrapExtError(err)}
			close(outputChan)
			return
		}
		defer rows.Close()
		for rows.Next() {
			item := ModelItems.DataChargesByUserItem{}
			err = db.SqlScanStruct(rows, &item)
			if err != nil {
				outputChan <- ChanChargesSearch{nil, exterror.WrapExtError(err)}
				close(outputChan)
				return
			}
			item.SearchDateDetail = condition["search_date"].(string)
			outputChan <- ChanChargesSearch{&item, nil}
		}
		close(outputChan)
	}()
	return outputChan
}

// Insert record after search report from BQ
func (this *BQLogSearchChargingModel) InsertLogChargingByUser(itemInsert ModelItems.BQLogSearchChargingItem) error {

	args := []interface{}{}
	args = append(args, itemInsert.ServerName)
	args = append(args, itemInsert.ShopCd)
	args = append(args, itemInsert.ShopName)
	args = append(args, itemInsert.FranchiseCd)
	args = append(args, itemInsert.UserID)
	args = append(args, itemInsert.UserName)
	args = append(args, itemInsert.UserMenu)
	args = append(args, itemInsert.UseTBLSize)
	args = append(args, itemInsert.GCPCharging)
	args = append(args, itemInsert.VJCharging)
	args = append(args, itemInsert.ExecTime)
	args = append(args, itemInsert.AppVersion)
	args = append(args, itemInsert.SearchingDate)
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD START
	args = append(args, itemInsert.Handle)
	args = append(args, itemInsert.Format)
	args = append(args, itemInsert.Tab)
	args = append(args, itemInsert.AppID)
	// ASO-5679 [BA]mBAWEB-v11a 初速比較：グラフの色の変更 - ADD END

	sqlStrings := `
INSERT INTO
	bq_log_search_charging (
		bqls_create_date,
		bqls_server_name,
		bqls_shop_cd,
		bqls_shop_name,
		bqls_franchise_cd,
		bqls_user_ID,
		bqls_user_name,
		bqls_use_menu,
		bqls_use_TBL_size,
		bqls_GCP_charging,
		bqls_VJ_charging,
		bqls_exec_time,
		bqls_app_version,
		bqls_searching_date,
		bqls_app_handle,
		bqls_app_format,
		bqls_app_tab,
		bqls_app_id
) VALUES (now(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	_, err := this.DB.Exec(sqlStrings, args...)
	if err != nil {
		return exterror.WrapExtError(err)
	}

	return nil
}
