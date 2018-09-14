package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
)

type CustomerModel struct {
	DB *sql.DB
}

func (this *CustomerModel) GetCustomer(userId string) ([]ModelItems.Customer, error) {
	um := UserMasterModel{this.DB}

	role, err := um.GetUserRole(userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	if role == nil {
		return nil, nil
	}

	listCustomer := []ModelItems.Customer{}

	sqlVar := []interface{}{}
	var userShopCondition string

	if !role.IsHonbu {
		userShopCondition = ` AND shm_shop_cd = ? `
		sqlVar = append(sqlVar, role.ShopCD)
	} else {
		if role.ChainCD != "" && len(role.ListChainCD) > 0 {
			userShopCondition = ` AND shm_chain_cd IN (` + Common.SQLPara(role.ListChainCD) + `) `
			sqlVar = append(sqlVar, Common.ToInterfaceArray(role.ListChainCD)...)
		} else if role.FranchiseCd != "" {
			userShopCondition = ` AND shm_franchise_cd = ? `
			sqlVar = append(sqlVar, role.FranchiseCd)
		} else if role.ShopCD != "" {
			userShopCondition = ` AND shm_shop_cd = ? `
			sqlVar = append(sqlVar, role.ShopCD)
		} else {
			userShopCondition = ` AND FALSE `
		}
	}

	sqlString := `
SELECT
	mc_server_name,
    mc_shop_cd,
    mc_customer_cat_cd,
    mc_customer_cat_name,
    mc_gender_cd,
    mc_gender_name,
    mc_flg_del
FROM master_customer
WHERE
     mc_shop_cd IN (
        SELECT shm_shop_cd
        FROM shop_master
        WHERE shm_flg_ope = '0' ` + userShopCondition + `
    )
    AND mc_flg_del = '0'
`
	rows, err := this.DB.Query(sqlString, sqlVar...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newItem := ModelItems.Customer{}
		err = db.SqlScanStruct(rows, &newItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listCustomer = append(listCustomer, newItem)
	}

	return listCustomer, nil
}

func (this *CustomerModel) GetCustomerCat() (*[]ModelItems.Customer, error) {
	results := []ModelItems.Customer{}
	sqlString := `
SELECT DISTINCT
	mc_customer_cat_name
FROM
	master_customer
ORDER BY
	mc_customer_cat_cd
`
	rows, err := this.DB.Query(sqlString)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newItem := ModelItems.Customer{}
		err = db.SqlScanStruct(rows, &newItem)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		results = append(results, newItem)
	}
	return &results, nil
}
