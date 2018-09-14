package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"errors"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"strings"
)

type ShopMasterModel struct {
	DB *sql.DB
}

func (this *ShopMasterModel) GetListShopByUser(userId string) ([]ModelItems.ShopItem, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listShop := []ModelItems.ShopItem{}

	sqlVar := []interface{}{}

	sqlString := ``
	if !role.IsHonbu {
		sqlString = `
SELECT
	shm_server_name,
    shm_shop_cd,
    shm_shop_name,
    shm_shared_book_store_code
FROM
    shop_master_show
WHERE
	shm_server_name = ?
	AND shm_shop_cd = ?
ORDER BY
	shm_chain_cd,
	shm_shop_seq_number,
	shm_shop_cd
`
		sqlVar = append(sqlVar, role.ServerName, role.ShopCD)
	} else {
		if role.FranchiseGroupCd != "" {
			sqlString = `
SELECT
	shm_server_name,
    shm_shop_cd,
    shm_shop_name,
    shm_shared_book_store_code
FROM
	(SELECT
		sfgm_franchise_cd
	FROM shop_franchise_group_master
	WHERE
		sfgm_delete_flag = '0'
		AND sfgm_franchise_group_cd = ?
	) fcd
JOIN shop_master_show shm
	ON sfgm_franchise_cd = shm_franchise_cd
ORDER BY
	shm.shm_chain_cd,
	shm.shm_shop_seq_number,
	shm.shm_shop_cd
`
			sqlVar = append(sqlVar, role.FranchiseGroupCd)
		} else if role.FranchiseCd != "" {
			sqlString = `
SELECT
	shm_server_name,
    shm_shop_cd,
    shm_shop_name,
    shm_shared_book_store_code
FROM
    shop_master_show
WHERE
	shm_franchise_cd = ?
ORDER BY
	shm_chain_cd,
	shm_shop_seq_number,
	shm_shop_cd
	`
			sqlVar = append(sqlVar, role.FranchiseCd)
		} else if role.ShopCD != "" {
			sqlString = `
SELECT
	shm_server_name,
    shm_shop_cd,
    shm_shop_name,
    shm_shared_book_store_code
FROM
    shop_master_show
WHERE
	shm_server_name = ?
	AND shm_shop_cd = ?
ORDER BY
	shm_chain_cd,
	shm_shop_seq_number,
	shm_shop_cd
`
			sqlVar = append(sqlVar, role.ServerName, role.ShopCD)
		} else {
			sqlString += `
SELECT
	shm_server_name,
    shm_shop_cd,
    shm_shop_name,
    shm_shared_book_store_code
FROM
    shop_master_show
WHERE FALSE
ORDER BY
	shm_chain_cd,
	shm_shop_seq_number,
	shm_shop_cd
`
		}
	}

	rows, err := this.DB.Query(sqlString, sqlVar...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newShopItem := ModelItems.ShopItem{}
		err = db.SqlScanStruct(rows, &newShopItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listShop = append(listShop, newShopItem)
	}
	return listShop, nil
}

func (this *ShopMasterModel) GetListShopByShopCode(arrShop []string) (*[]ModelItems.ShopItem, error) {
	if len(arrShop) == 0 {
		return nil, exterror.WrapExtError(errors.New("arrShop is empty"))
	}

	strCondition := " shm_shop_cd "
	if len(arrShop) > 0 && strings.Contains(arrShop[0], "|") {
		strCondition = " CONCAT(shm_server_name,'|',shm_shop_cd) "
	}
	strCondition += ` IN (?` + strings.Repeat(",?", len(arrShop)-1) + `)`
	sqlString := `
SELECT
	shm_shop_cd,
	shm_shop_name,
	shm_server_name,
	shm_chain_cd,
	shm_franchise_cd
FROM
	shop_master_show
WHERE ` + strCondition
	listShop := []ModelItems.ShopItem{}
	rows, err := this.DB.Query(sqlString, Common.ToInterfaceArray(arrShop)...)
	if err != nil {
		return &listShop, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newItem := ModelItems.ShopItem{}
		err := db.SqlScanStruct(rows, &newItem)
		if err != nil {
			return &listShop, exterror.WrapExtError(err)
		}
		listShop = append(listShop, newItem)
	}
	return &listShop, nil
}

func (this *ShopMasterModel) GetInfoShopByCD(shopCd string) (ModelItems.ShopItem, error) {
	shop := ModelItems.ShopItem{}
	sqlString := `
SELECT
	shm_shop_cd,
	shm_shop_name,
	shm_server_name,
	shm_chain_cd,
	shm_franchise_cd,
	shm_shared_book_store_code,
	shm_tel_no,
	shm_biz_start_time,
	shm_biz_end_time,
	shm_address,
	shm_shop_name_2
FROM
	shop_master_show
WHERE shm_shop_cd = ?
`
	rows, err := this.DB.Query(sqlString, shopCd)
	if err != nil {
		return shop, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := db.SqlScanStruct(rows, &shop)
		if err != nil {
			return shop, exterror.WrapExtError(err)
		}
	}
	return shop, nil
}
