package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"strings"
)

type BQCategoryModel struct {
	DB *sql.DB
}

//BSPOS分類
func (this *BQCategoryModel) GetBumonListByUser(userId string) ([]ModelItems.BQCategoryCDNameItem, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listCategoryItem := []ModelItems.BQCategoryCDNameItem{}

	sqlVar := []interface{}{}
	var userShopCondition string
	if strings.Contains(role.ChainCD, Common.CHAIN_VJ) {
		userShopCondition += ` AND TRUE `
	} else {
		if !role.IsHonbu {
			userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
			sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
		} else {
			if role.FranchiseGroupCd != "" {
				userShopCondition += `
AND shm_franchise_cd IN (
	SELECT sfgm_franchise_cd
	FROM shop_franchise_group_master
	WHERE sfgm_franchise_group_cd = ?)
`
				sqlVar = append(sqlVar, role.FranchiseGroupCd)
			} else if role.FranchiseCd != "" {
				userShopCondition += `
                AND shm_franchise_cd = ?`
				sqlVar = append(sqlVar, role.FranchiseCd)
			} else {
				userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
				sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
			}
		}
	}

	sqlString := `
SELECT
    bqct_bumon_cd CD,
    bqct_bumon_nm Name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '3'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_bumon_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_bumon_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_bumon_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_bumon_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_bumon_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_bumon_cd,",%")
            ))
    ` + userShopCondition + `
GROUP BY
    bqct_bumon_cd, bqct_bumon_nm
`
	rows, err := this.DB.Query(sqlString, sqlVar...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCategoryItem := ModelItems.BQCategoryCDNameItem{}
		err = db.SqlScanStruct(rows, &newCategoryItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listCategoryItem = append(listCategoryItem, newCategoryItem)
	}

	return listCategoryItem, nil
}

//店舗POS分類
func (this *BQCategoryModel) GetKubunListByUser(userId string) ([]ModelItems.BQCategoryCDNameItem, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listCategoryItem := []ModelItems.BQCategoryCDNameItem{}

	sqlVar := []interface{}{}
	var userShopCondition string
	if strings.Contains(role.ChainCD, Common.CHAIN_VJ) {
		userShopCondition += ` AND TRUE `
	} else {
		if !role.IsHonbu {
			userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
			sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
		} else {
			if role.FranchiseGroupCd != "" {
				userShopCondition += `
AND shm_franchise_cd IN (
	SELECT sfgm_franchise_cd
	FROM shop_franchise_group_master
	WHERE sfgm_franchise_group_cd = ?)
`
				sqlVar = append(sqlVar, role.FranchiseGroupCd)
			} else if role.FranchiseCd != "" {
				userShopCondition += `
                AND shm_franchise_cd = ?`
				sqlVar = append(sqlVar, role.FranchiseCd)
			} else {
				userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
				sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
			}
		}
	}

	sqlString := `
SELECT
    bqct_kubn_cd CD,
    bqct_kubn_nm Name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '4'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_kubn_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_kubn_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_kubn_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_kubn_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_kubn_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_kubn_cd,",%")
            ))
    ` + userShopCondition + `
GROUP BY
    bqct_kubn_cd, bqct_kubn_nm
`
	rows, err := this.DB.Query(sqlString, sqlVar...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCategoryItem := ModelItems.BQCategoryCDNameItem{}
		err = db.SqlScanStruct(rows, &newCategoryItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listCategoryItem = append(listCategoryItem, newCategoryItem)
	}

	return listCategoryItem, nil
}

//MSグループ分類
func (this *BQCategoryModel) GetMediaGroupListByUser(userId string) ([]ModelItems.BQCategoryCDNameItem, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listCategoryItem := []ModelItems.BQCategoryCDNameItem{}

	sqlVar := []interface{}{}
	var userShopCondition string
	if strings.Contains(role.ChainCD, Common.CHAIN_VJ) {
		userShopCondition += ` AND TRUE `
	} else {
		if !role.IsHonbu {
			userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
			sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
		} else {
			if role.FranchiseGroupCd != "" {
				userShopCondition += `
AND shm_franchise_cd IN (
	SELECT sfgm_franchise_cd
	FROM shop_franchise_group_master
	WHERE sfgm_franchise_group_cd = ?)
`
				sqlVar = append(sqlVar, role.FranchiseGroupCd)
			} else if role.FranchiseCd != "" {
				userShopCondition += `
                AND shm_franchise_cd = ?`
				sqlVar = append(sqlVar, role.FranchiseCd)
			} else {
				userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
				sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
			}
		}
	}

	sqlString := `
SELECT
    bccm_media_group_cd CD,
    bccm_media_group_name Name
FROM bq_category_ms
GROUP BY
    bccm_media_group_cd, bccm_media_group_name
`
	rows, err := this.DB.Query(sqlString)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCategoryItem := ModelItems.BQCategoryCDNameItem{}
		err = db.SqlScanStruct(rows, &newCategoryItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listCategoryItem = append(listCategoryItem, newCategoryItem)
	}

	return listCategoryItem, nil

}

//MS大中小メディア分類
func (this *BQCategoryModel) GetMediaListByUser(userId string) ([]ModelItems.BQMediaItem, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listCategoryItem := []ModelItems.BQMediaItem{}

	sqlVar := []interface{}{}
	var userShopCondition string
	if strings.Contains(role.ChainCD, Common.CHAIN_VJ) {
		userShopCondition += ` AND TRUE `
	} else {
		if !role.IsHonbu {
			userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
			sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
		} else {
			if role.FranchiseGroupCd != "" {
				userShopCondition += `
AND shm_franchise_cd IN (
	SELECT sfgm_franchise_cd
	FROM shop_franchise_group_master
	WHERE sfgm_franchise_group_cd = ?)
`
				sqlVar = append(sqlVar, role.FranchiseGroupCd)
			} else if role.FranchiseCd != "" {
				userShopCondition += `
                AND shm_franchise_cd = ?`
				sqlVar = append(sqlVar, role.FranchiseCd)
			} else {
				userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
				sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
			}
		}
	}

	sqlString := `
SELECT *
FROM (
(SELECT DISTINCT
	'1' media_type,
    bccm_media_group1_cd media_cd,
    bccm_media_group1_name media_name
FROM bq_category_ms)

UNION ALL

(SELECT DISTINCT
	'2' media_type,
    bccm_media_group2_cd media_cd,
    bccm_media_group2_name media_name
FROM bq_category_ms)

UNION ALL

(SELECT DISTINCT
	'3' media_type,
    bccm_media_group3_cd media_cd,
    bccm_media_group3_name media_name
FROM bq_category_ms)

UNION ALL

(SELECT DISTINCT
	'4' media_type,
    bccm_media_group4_cd media_cd,
    bccm_media_group4_name media_name
FROM bq_category_ms)
) tmp
ORDER BY media_type, media_cd
`
	sqlVarGroup := []interface{}{}
	sqlVarGroup = append(sqlVarGroup, sqlVar...)

	rows, err := this.DB.Query(sqlString)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCategoryItem := ModelItems.BQMediaItem{}
		err = db.SqlScanStruct(rows, &newCategoryItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listCategoryItem = append(listCategoryItem, newCategoryItem)
	}

	return listCategoryItem, nil

}

//MS中メディア分類
func (this *BQCategoryModel) GetMediaGroup2ListByUser(userId string) ([]ModelItems.BQCategoryCDNameItem, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listCategoryItem := []ModelItems.BQCategoryCDNameItem{}

	sqlVar := []interface{}{}
	var userShopCondition string
	if strings.Contains(role.ChainCD, Common.CHAIN_VJ) {
		userShopCondition += ` AND TRUE `
	} else {
		if !role.IsHonbu {
			userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
			sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
		} else {
			if role.FranchiseGroupCd != "" {
				userShopCondition += `
AND shm_franchise_cd IN (
	SELECT sfgm_franchise_cd
	FROM shop_franchise_group_master
	WHERE sfgm_franchise_group_cd = ?)
`
				sqlVar = append(sqlVar, role.FranchiseGroupCd)
			} else if role.FranchiseCd != "" {
				userShopCondition += `
                AND shm_franchise_cd = ?`
				sqlVar = append(sqlVar, role.FranchiseCd)
			} else {
				userShopCondition += `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
				sqlVar = append(sqlVar, role.ShopCD, role.ServerName)
			}
		}
	}

	sqlString := `
SELECT
    bqct_media_group2_cd CD,
    bqct_media_group2_name Name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '2'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
	bqct_media_group2_cd != ''
    AND (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_media_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_media_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_media_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_media_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd,",%")
            ))
    ` + userShopCondition + `
GROUP BY
    CD, Name
`
	sqlVarGroup := []interface{}{}
	sqlVarGroup = append(sqlVarGroup, sqlVar...)

	rows, err := this.DB.Query(sqlString, sqlVarGroup...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCategoryItem := ModelItems.BQCategoryCDNameItem{}
		err = db.SqlScanStruct(rows, &newCategoryItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listCategoryItem = append(listCategoryItem, newCategoryItem)
	}

	return listCategoryItem, nil

}

func (this *BQCategoryModel) GetGenreListByShopCdAndFranchiseCd(franchiseCd, serverName, shopCd string, listBumonItem *[]ModelItems.BQCategoryCDNameItem, listKubunItem *[]ModelItems.BQCategoryCDNameItem, listMediaItem *[]ModelItems.BQMediaItem) (string, string, error) {
	productCd := ""
	referGenreType := ""
	query := `
SELECT shm_product_cd,ssr_refer_genre_type
FROM shop_master_show
JOIN setting_shop_refer
ON ssr_server_name = shm_server_name
AND ssr_shop_cd = shm_shop_cd
WHERE shm_franchise_cd = ?
AND shm_server_name = ?
AND shm_shop_cd = ?
	`
	err := this.DB.QueryRow(query, franchiseCd, serverName, shopCd).Scan(&productCd, &referGenreType)
	if err != nil {
		return productCd, referGenreType, exterror.WrapExtError(err)
	}
	sqlVar := []interface{}{}
	condition := " AND shm_franchise_cd = ? AND shm_server_name = ? AND shm_shop_cd = ?"
	sqlVar = append(sqlVar, franchiseCd, serverName, shopCd)
	if productCd == "M" {
		query = `
(SELECT
	'1' media_type,
    bqct_media_group1_cd media_cd,
    bqct_media_group1_name media_name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '2'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_media_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_media_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_media_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_media_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd,",%")
            ))
    ` + condition + `
GROUP BY
    media_cd, media_name)

UNION ALL

(SELECT
	'2' media_type,
    bqct_media_group2_cd media_cd,
    bqct_media_group2_name media_name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '2'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_media_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_media_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_media_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_media_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd,",%")
            ))
    ` + condition + `
GROUP BY
    media_cd, media_name)

UNION ALL

(SELECT
	'3' media_type,
    bqct_media_group3_cd media_cd,
    bqct_media_group3_name media_name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '2'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_media_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_media_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_media_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_media_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd,",%")
            ))
    ` + condition + `
GROUP BY
    media_cd, media_name)

UNION ALL

(SELECT
	'4' media_type,
    bqct_media_group4_cd media_cd,
    bqct_media_group4_name media_name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '2'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_media_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_media_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_media_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_media_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_media_cd,",%")
            ))
    ` + condition + `
GROUP BY
    media_cd, media_name)
		`
		sqlVarGroup := []interface{}{}
		sqlVarGroup = append(sqlVarGroup, sqlVar...)
		sqlVarGroup = append(sqlVarGroup, sqlVar...)
		sqlVarGroup = append(sqlVarGroup, sqlVar...)
		sqlVarGroup = append(sqlVarGroup, sqlVar...)

		rows, err := this.DB.Query(query, sqlVarGroup...)
		if err != nil {
			return productCd, referGenreType, exterror.WrapExtError(err)
		}
		defer rows.Close()
		for rows.Next() {
			newCategoryItem := ModelItems.BQMediaItem{}
			err = db.SqlScanStruct(rows, &newCategoryItem)
			if err != nil {
				return productCd, referGenreType, exterror.WrapExtError(err)
			}
			*listMediaItem = append(*listMediaItem, newCategoryItem)
		}
	} else if productCd == "B" {
		if referGenreType == "3" {
			query = `
		SELECT
    bqct_bumon_cd CD,
    bqct_bumon_nm Name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '3'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_bumon_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_bumon_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_bumon_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_bumon_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_bumon_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_bumon_cd,",%")
            ))
    ` + condition + `
GROUP BY
    bqct_bumon_cd, bqct_bumon_nm
`
			rows, err := this.DB.Query(query, sqlVar...)
			if err != nil {
				return productCd, referGenreType, exterror.WrapExtError(err)
			}
			defer rows.Close()
			for rows.Next() {
				newCategoryItem := ModelItems.BQCategoryCDNameItem{}
				err = db.SqlScanStruct(rows, &newCategoryItem)
				if err != nil {
					return productCd, referGenreType, exterror.WrapExtError(err)
				}
				*listBumonItem = append(*listBumonItem, newCategoryItem)
			}

		} else if referGenreType == "4" {

			query = `
SELECT
    bqct_kubn_cd CD,
    bqct_kubn_nm Name
FROM bq_category
JOIN setting_shop_refer
ON
    bqct_shop_cd = ssr_shop_cd
    AND bqct_servername = ssr_server_name
    AND ssr_refer_genre_type = '4'
JOIN shop_master_show
ON
	bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
WHERE
    (IFNULL(ssr_start_genre_cd, '') = '' OR bqct_kubn_cd >= ssr_start_genre_cd)
    AND (IFNULL(ssr_end_genre_cd, '') = '' OR bqct_kubn_cd <= ssr_end_genre_cd)
    AND (IFNULL(ssr_genre_cd_check_off, '') = '' OR NOT(
                ssr_genre_cd_check_off LIKE bqct_kubn_cd
                OR ssr_genre_cd_check_off LIKE CONCAT(bqct_kubn_cd,",%")
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_kubn_cd)
                OR ssr_genre_cd_check_off LIKE CONCAT("%,",bqct_kubn_cd,",%")
            ))
    ` + condition + `
GROUP BY
    bqct_kubn_cd, bqct_kubn_nm
`
			rows, err := this.DB.Query(query, sqlVar...)
			if err != nil {
				return productCd, referGenreType, exterror.WrapExtError(err)
			}
			defer rows.Close()
			for rows.Next() {
				newCategoryItem := ModelItems.BQCategoryCDNameItem{}
				err = db.SqlScanStruct(rows, &newCategoryItem)
				if err != nil {
					return productCd, referGenreType, exterror.WrapExtError(err)
				}
				*listKubunItem = append(*listKubunItem, newCategoryItem)
			}

		}

	}

	return productCd, referGenreType, nil
}

func (this *BQCategoryModel) MediaListByRoleUser(userId string) (*[]ModelItems.MediaItem, error) {

	listCategoryItem := []ModelItems.MediaItem{}

	um := UserMasterModel{this.DB}
	role, err := um.GetUserRole(userId)
	if err != nil {
		return &listCategoryItem, exterror.WrapExtError(err)
	}
	if role == nil {
		return &listCategoryItem, nil
	}

	args := []interface{}{}
	condition := `
AND shm_franchise_cd IN (
	SELECT sfgm_franchise_cd
	FROM shop_franchise_group_master
	WHERE ? = '' OR sfgm_franchise_group_cd = ?
)
`
	if strings.Contains(role.ChainCD, Common.CHAIN_VJ) {
		args = append(args, "", "")
	} else {
		if !role.IsHonbu {
			condition = `
AND shm_shop_cd = ?
AND shm_server_name = ?
`
			args = append(args, role.ShopCD, role.ServerName)
		} else {
			if role.FranchiseGroupCd != "" {
				args = append(args, role.FranchiseGroupCd, role.FranchiseGroupCd)
			} else if role.FranchiseCd != "" {
				condition = `
AND shm_franchise_cd  = ?
`
				args = append(args, role.FranchiseCd)
			}
		}
	}

	query := `
SELECT DISTINCT
	bqct_media_cd media_cd,
	bqct_media_name media_name,
	shm_product_cd product
FROM
(
	SELECT
		DISTINCT
		bqct_shop_cd, bqct_servername,bqct_category_type,shm_product_cd,
		CASE
			WHEN bqct_category_type IN('1','2') THEN IFNULL(bqct_media_group2_cd, '')
			WHEN bqct_category_type IN('3','4') THEN IFNULL(bqct_bumon_cd, '')
		END bqct_media_cd
		,
		CASE
			WHEN bqct_category_type IN('1','2') THEN bqct_media_group2_name
			WHEN bqct_category_type IN('3','4') THEN bqct_bumon_nm
		END bqct_media_name
	FROM bq_category
	JOIN shop_master_show
	  ON bqct_shop_cd = shm_shop_cd
    AND bqct_servername = shm_server_name
    ` + condition + `
) bqct
JOIN setting_shop_refer ssr
	ON bqct.bqct_shop_cd = ssr.ssr_shop_cd
	AND bqct.bqct_servername = ssr.ssr_server_name
	AND bqct.bqct_category_type = ssr.ssr_refer_genre_type
	AND ( ssr.ssr_start_genre_cd = '' OR bqct.bqct_media_cd >= ssr.ssr_start_genre_cd )
	AND ( ssr.ssr_end_genre_cd = '' OR bqct.bqct_media_cd <= ssr.ssr_end_genre_cd )
	AND ssr.ssr_genre_cd_check_off NOT LIKE bqct.bqct_media_cd
	AND ssr.ssr_genre_cd_check_off NOT LIKE CONCAT('%,',bqct.bqct_media_cd)
	AND ssr.ssr_genre_cd_check_off NOT LIKE CONCAT(bqct.bqct_media_cd,',%')
	AND ssr.ssr_genre_cd_check_off NOT LIKE CONCAT('%,',bqct.bqct_media_cd,',%')
ORDER BY bqct_media_cd
`
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return &listCategoryItem, exterror.WrapExtError(err)
	}
	defer rows.Close()

	for rows.Next() {
		newCategoryItem := ModelItems.MediaItem{}
		err = gf.SqlScanStruct(rows, &newCategoryItem)
		if err != nil {
			return &listCategoryItem, exterror.WrapExtError(err)
		}
		listCategoryItem = append(listCategoryItem, newCategoryItem)
	}

	return &listCategoryItem, nil
}

// MSメディア大分類
func (this *BQCategoryModel) GetMediaGroup1List() ([]ModelItems.MediaItem, error) {

	listCategoryItem := []ModelItems.MediaItem{}

	sqlString := `
SELECT
    bccm_media_group1_cd media_cd,
    bccm_media_group1_name media_name
FROM bq_category_ms
WHERE
	bccm_media_group1_cd != ''
GROUP BY
    media_cd, media_name
`

	rows, err := this.DB.Query(sqlString)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newCategoryItem := ModelItems.MediaItem{}
		err = db.SqlScanStruct(rows, &newCategoryItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listCategoryItem = append(listCategoryItem, newCategoryItem)
	}

	return listCategoryItem, nil

}


func (this *BQCategoryModel) ListFullGenreByJanGroup(janGroup []string) ([]*ModelItems.MediaGroup1,error) {
	var listMedia []*ModelItems.MediaGroup1

	args := []interface{}{}

	for _,v := range janGroup {
		args = append(args,v)
	}

	sWhereJanGroup := ""
	if len(janGroup) > 0 {
		sWhereJanGroup = ` AND bcslg.jan_grouping IN (?` + strings.Repeat(",?", len(janGroup)-1) + `)`
	}

	sql := `
SELECT
	bccm.bccm_media_group1_cd,
	bccm.bccm_media_group1_name,
	bccm.bccm_media_group2_cd,
	bccm.bccm_media_group2_name,
	bccm.bccm_media_group3_cd,
	bccm.bccm_media_group3_name,
	bccm.bccm_media_group4_cd,
	bccm.bccm_media_group4_name,
	bcslg.jan_grouping
FROM  bq_category_super_large_grouping bcslg
LEFT JOIN bq_category_ms bccm
ON bcslg.bccm_media_group1_cd = bccm.bccm_media_group1_cd
WHERE TRUE
` + sWhereJanGroup + `
ORDER BY
	bccm.bccm_media_group1_cd,
	bccm.bccm_media_group2_cd,
	bccm.bccm_media_group3_cd,
	bccm.bccm_media_group4_cd
	`
	rows, err := this.DB.Query(sql, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	mediaCD1 := ""
	mediaCD2 := ""
	mediaCD3 := ""
	mediaCD4 := ""
	mediaName1 := ""
	mediaName2 := ""
	mediaName3 := ""
	mediaName4 := ""
	janGroupCD := ""

	mapMediaGroup1 := map[string]*ModelItems.MediaGroup1{}
	mapMediaGroup2 := map[string]*ModelItems.MediaGroup2{}
	mapMediaGroup3 := map[string]*ModelItems.MediaGroup3{}

	for rows.Next() {
		err = rows.Scan(&mediaCD1, &mediaName1, &mediaCD2, &mediaName2, &mediaCD3, &mediaName3, &mediaCD4, &mediaName4, &janGroupCD)

		if mapMediaGroup1[mediaCD1] == nil {
			newMG1 := ModelItems.MediaGroup1{}
			newMG1.MediaGroup1Cd = mediaCD1
			newMG1.MediaGroup1Name = mediaName1
			newMG1.JanGroup = janGroupCD
			newMG1.MediaGroup2 = []*ModelItems.MediaGroup2{}
			listMedia = append(listMedia,&newMG1)
			mapMediaGroup1[mediaCD1] = &newMG1
		}
		if mapMediaGroup2[mediaCD1+mediaCD2] == nil && mediaCD2 != "" {
			newMG2 := ModelItems.MediaGroup2{}
			newMG2.MediaGroup2Cd = mediaCD2
			newMG2.MediaGroup2Name = mediaName2
			newMG2.MediaGroup3 = []*ModelItems.MediaGroup3{}
			mapMediaGroup1[mediaCD1].MediaGroup2 = append(mapMediaGroup1[mediaCD1].MediaGroup2, &newMG2)
			mapMediaGroup2[mediaCD1 + mediaCD2] = &newMG2
		}
		if mapMediaGroup3[mediaCD1+mediaCD2+mediaCD3] == nil && mediaCD3 != "" {
			newMG3 := ModelItems.MediaGroup3{}
			newMG3.MediaGroup3Cd = mediaCD3
			newMG3.MediaGroup3Name = mediaName3
			newMG3.MediaGroup4 = []*ModelItems.MediaGroup4{}
			mapMediaGroup2[mediaCD1+mediaCD2].MediaGroup3 = append(mapMediaGroup2[mediaCD1+mediaCD2].MediaGroup3, &newMG3)
			mapMediaGroup3[mediaCD1 + mediaCD2 +mediaCD3] = &newMG3
		}
		if mediaCD4 != "" {
			mapMediaGroup3[mediaCD1 + mediaCD2 + mediaCD3].MediaGroup4 = append(mapMediaGroup3[mediaCD1+mediaCD2+mediaCD3].MediaGroup4, &ModelItems.MediaGroup4{mediaCD4,mediaName4})
		}

	}


	return listMedia,nil
}