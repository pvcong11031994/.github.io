package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"errors"
	"fmt"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"strings"
)

const (
	flag_auto   = "auto"
	flag_manual = "manual"

	flag_edi_data   = "data"
	flag_edi_nodata = "nodata"
)

type WorkTransOrderListFixModel struct {
	DB *sql.DB
}

type WtOrderData struct {
	CreateDate              string `sql:"wtolf_create_date"`
	UpdateDate              string `sql:"wtolf_update_date"`
	ReceiptNo               int    `sql:"wtolf_receipt_no"`
	ShopCd                  string `sql:"wtolf_shop_cd"`
	Jan                     string `sql:"wtolf_jan"`
	OrderCount              int    `sql:"wtolf_order_count"`
	OrderWayCd              string `sql:"wtolf_order_way_cd"`
	FlgEdi                  string `sql:"wtolf_flg_edi"`
	OrderSupplierCd         string `sql:"wtolf_order_supplier_cd"`
	OrderShopLineNo         string `sql:"wtolf_order_shop_line_no"`
	OrderShopCd             string `sql:"wtolf_order_shop_cd"`
	OrderMemo               string `sql:"wtolf_order_memo"`
	SupplierOrderCd         string `sql:"wtolf_supplier_order_cd"`
	FlgItemNew              string `sql:"wtolf_flg_item_new"`
	FlgAutoOrder            string `sql:"wtolf_flg_auto_order"`
	FlgStopOrder            string `sql:"wtolf_flg_stop_order"`
	FlgOrderControl         string `sql:"wtolf_flg_order_control"`
	FlgCustomerOrderControl string `sql:"wtolf_flg_customer_order_control"`
	FlgItemGroupControl     string `sql:"wtolf_flg_item_group_control"`
	FlgUnexpectedValue      string `sql:"wtolf_flg_unexpected_value"`
	FlgOrderCalcWay         string `sql:"wtolf_flg_order_calc_way"`
	FlgDataFixDate          string `sql:"wtolf_flg_data_fix_date"`
	FlgOrderExport          string `sql:"wtolf_flg_order_export"`
	OpeStatusFls            string `sql:"wtolf_ope_status_fls"`
	DateFixDate             string `sql:"wtolf_date_fix_date"`
	PublisherCd             string `sql:"wtolf_publisher_cd"`
	PublisherName           string `sql:"wtolf_publisher_name"`
	GenreCd                 string `sql:"wtolf_genre_cd"`
	GenreName               string `sql:"wtolf_genre_name"`
	ReceiptNoLb             string `sql:"wtolf_receipt_no_lb"`
	BgenreCd                string `sql:"wtolf_Bgenre_cd"`
	LocationCd              string `sql:"wtolf_location_cd"`
	ServerName              string `sql:"wtolf_server_name"`
	FranchiseCd             string `sql:"wtolf_franchise_cd"`
}

type ConfirmOrderData struct {
	GmJan         string `sql:"gm_jan"`
	GoodsName     string `sql:"gm_goods_name"`
	Author        string `sql:"gm_author"`
	PublisherName string `sql:"gm_publisher_name"`
	PublisherCd   string `sql:"gm_publisher_cd"`
	BGenreCd      string `sql:"gm_media_cd"`
	LocationCd    string `sql:"location_cd"`

	FranchiseCd string `sql:"wtolf_franchise_cd"`
	ServerName  string `sql:"wtolf_server_name"`
	ShopCd      string `sql:"wtolf_shop_cd"`
	ShopName    string `sql:"shop_name"`

	Jan             string `sql:"wtolf_jan"`
	OrderCount      int    `sql:"wtolf_order_count"`
	FlgAutoOrder    string `sql:"wtolf_flg_auto_order"`
	FlgEdi          string `sql:"wtolf_flg_edi"`
	OrderSupplierCd string `sql:"wtolf_order_supplier_cd"`
	OrderShopLineNo string `sql:"wtolf_order_shop_line_no"`
	OrderShopCd     string `sql:"wtolf_order_shop_cd"`

	BqctJan   string `sql:"bqct_jan_cd"`
	GenreCd   string `sql:"genre_cd"`
	GenreName string `sql:"genre_name"`
}

type WtOrderListData struct {
	OrderDate       string `sql:"wtolf_flg_data_fix_date" header:"発注日"`
	SupplierName    string `sql:"mw_supplier_name" header:"発注先名"`
	OrderShopLineNo string `sql:"wtolf_order_shop_line_no" header:"番線"`
	OrderShopCd     string `sql:"wtolf_order_shop_cd" header:"書店コード"`
	OrderMemo       string `sql:"wtolf_order_memo" header:"発注メモ"`
	ISBN            string `sql:"wtolf_jan" header:"ISBN"`
	GoodsName       string `sql:"gm_goods_name" header:"書名"`
	MakerName       string `sql:"wtolf_publisher_name" header:"出版社名"`
	Author          string `sql:"gm_author" header:"著者名"`
	Price           int    `sql:"gm_price" header:"価格"`
	OrderCount      int    `sql:"wtolf_order_count" header:"発注数"`
	ShopCD          string `sql:"wtolf_shop_cd" header:"店舗コード"`
	ShopName        string `sql:"shm_shop_name" header:"店舗名"`
	GenreCd         string `sql:"wtolf_genre_cd" header:"部門コード"`
	GenreName       string `sql:"genre_name" header:"ジャンル名"`
	EdiFlag         string `sql:"wtolf_flg_edi" header:"EDI区分名"`
	OrderWayName    string `sql:"mow_order_way_name" header:"発注方法"`
}

type WtOrderListDataWithErr struct {
	Data *WtOrderListData
	Err  error
}

func (this *WorkTransOrderListFixModel) SingleShopConfirm(shop ModelItems.ShopItem, jan_cd []string) (map[string]ConfirmOrderData, error) {
	listResult := map[string]ConfirmOrderData{}
	if len(jan_cd) == 0 {
		return listResult, nil
	}
	for _, jan := range jan_cd {
		listResult[jan] = ConfirmOrderData{
			ShopCd:      shop.ShopCD,
			ShopName:    shop.ShopName,
			ServerName:  shop.ServerName,
			FranchiseCd: shop.ShopFranchiseCD,
		}
	}
	query := `
SELECT
	gm_jan,
	gm_media_cd,
	gm_goods_name,
	gm_author,
	gm_publisher_name,
	gm_publisher_cd,
	location_cd
FROM
	ao_goods_master
LEFT JOIN
(
	SELECT
		bqls_jan_cd jan,
		bqls_tana_cd location_cd
	FROM
		bq_location_setting
	WHERE
	    bqls_tana_type = '1'
	AND bqls_server_name = ?
	AND bqls_shop_cd = ?
	AND bqls_jan_cd IN (` + Common.SQLPara(jan_cd) + `)
) bqls
	ON gm_jan = jan
WHERE
	gm_jan IN (` + Common.SQLPara(jan_cd) + `)
`
	args := []interface{}{shop.ServerName, shop.ShopCD}
	args = append(args, Common.ToInterfaceArray(jan_cd)...)
	args = append(args, Common.ToInterfaceArray(jan_cd)...)

	rows, err := this.DB.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	for rows.Next() {
		newItem := ConfirmOrderData{}
		err := db.SqlScanStruct(rows, &newItem)
		if err != nil {
			return listResult, exterror.WrapExtError(err)
		}
		newItem.Jan = newItem.GmJan
		listResult[newItem.Jan] = newItem
	}

	query = `
SELECT
	wtolf_jan,
	wtolf_order_count,
	wtolf_flg_auto_order,
	wtolf_server_name,
	wtolf_franchise_cd
FROM
	work_trans_order_list_fix
WHERE TRUE
	AND wtolf_jan IN (` + Common.SQLPara(jan_cd) + `)
	AND wtolf_shop_cd = ?
	AND wtolf_server_name = ?
	AND wtolf_franchise_cd = ?
	AND wtolf_create_date = ?
`
	argsWtolf := []interface{}{}
	argsWtolf = append(argsWtolf, Common.ToInterfaceArray(jan_cd)...)
	argsWtolf = append(argsWtolf, shop.ShopCD, shop.ServerName, shop.ShopFranchiseCD, Common.CurrentDate())
	rowsWtolf, err := this.DB.Query(query, argsWtolf...)

	defer rowsWtolf.Close()
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	for rowsWtolf.Next() {
		newItem := ConfirmOrderData{}
		err := db.SqlScanStruct(rowsWtolf, &newItem)
		if err != nil {
			return listResult, exterror.WrapExtError(err)
		}
		item := listResult[newItem.Jan]
		item.Jan = newItem.Jan
		item.OrderCount = newItem.OrderCount
		item.FlgAutoOrder = newItem.FlgAutoOrder
		listResult[newItem.Jan] = item
	}

	query = `
SELECT
	moa_server_name wtolf_server_name,
	moa_franchise_cd wtolf_franchise_cd,
	moa_shop_cd wtolf_shop_cd,
	mw_flg_edi wtolf_flg_edi,
	moa_supplier_cd wtolf_order_supplier_cd,
	moa_supplier_line_no wtolf_order_shop_line_no,
	moa_supplier_shop_cd wtolf_order_shop_cd
FROM
	ao_master_order_accounts
JOIN
	ao_master_wholeseller
 ON moa_supplier_cd = mw_supplier_cd
AND moa_franchise_cd = mw_franchise_cd
WHERE TRUE
	AND moa_shop_cd = ?
	AND moa_server_name = ?
	AND moa_franchise_cd = ?
`
	argsMoa := []interface{}{shop.ShopCD, shop.ServerName, shop.ShopFranchiseCD}
	rowsMoa, err := this.DB.Query(query, argsMoa...)
	defer rowsMoa.Close()
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	for rowsMoa.Next() {
		newItem := ConfirmOrderData{}
		err := db.SqlScanStruct(rowsMoa, &newItem)
		if err != nil {
			return listResult, exterror.WrapExtError(err)
		}
		for jan, item := range listResult {
			if item.ShopCd == newItem.ShopCd && item.FranchiseCd == newItem.FranchiseCd && item.ServerName == newItem.ServerName {
				item.FlgEdi = newItem.FlgEdi
				item.OrderSupplierCd = newItem.OrderSupplierCd
				item.OrderShopLineNo = newItem.OrderShopLineNo
				item.OrderShopCd = newItem.OrderShopCd
				listResult[jan] = item
			}
		}
	}

	query = `
SELECT
	bqct_servername,
	bqct_franchise_id,
	bqct_shop_cd,
	bqct_category_cd genre_cd,
	bqct_media_group_name genre_name,
	bqct_jan_cd
FROM bq_category
WHERE bqct_category_type IN ('1','2')
AND bqct_shop_cd = ?
AND bqct_servername = ?
UNION
SELECT
	bqct_servername,
	bqct_franchise_id,
	bqct_shop_cd,
	IF(ssr_refer_genre_type IN ('1','2') , bqct_kubn_cd, bqct_bumon_cd) genre_cd,
	IF(ssr_refer_genre_type IN ('1','2') , bqct_kubn_nm, bqct_bumon_nm) genre_name ,
	bqct_jan_cd
FROM (
SELECT
	bqct_servername,
	bqct_franchise_id,
	bqct_category_type,
	bqct_shop_cd,
	bqct_kubn_cd, bqct_kubn_nm,
	bqct_bumon_cd, bqct_bumon_nm,
	bqct_jan_cd
FROM bq_category
WHERE TRUE
	AND bqct_category_type IN ('3','4')
	AND bqct_shop_cd = ?
	AND bqct_servername = ?
) bqct
LEFT JOIN setting_shop_refer
	ON bqct_servername = ssr_server_name
	AND bqct_shop_cd = ssr_shop_cd
WHERE TRUE
	AND bqct_shop_cd = ?
	AND bqct_servername = ?
	AND bqct_jan_cd IN (` + strings.Trim(strings.Repeat("?,", len(jan_cd)), ",") + `)
`
	argsBqct := []interface{}{shop.ShopCD, shop.ServerName, shop.ShopCD, shop.ServerName, shop.ShopCD, shop.ServerName}
	argsBqct = append(argsBqct, Common.ToInterfaceArray(jan_cd)...)
	rowsBqct, err := this.DB.Query(query, argsBqct...)
	defer rowsBqct.Close()
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	for rowsBqct.Next() {
		newItem := ConfirmOrderData{}
		err := db.SqlScanStruct(rowsBqct, &newItem)
		if err != nil {
			return listResult, exterror.WrapExtError(err)
		}
		for jan, item := range listResult {
			if item.ShopCd == newItem.ShopCd && item.FranchiseCd == newItem.FranchiseCd && jan == newItem.BqctJan {
				item.GenreCd = newItem.GenreCd
				item.GenreName = newItem.GoodsName
				listResult[item.Jan] = item
			}
		}
	}
	return listResult, nil
}

/////////////
func (this *WorkTransOrderListFixModel) MultiShopConfirm(args, argsShopServer []interface{}, jan string) (map[string]ConfirmOrderData, error) {
	listResult := map[string]ConfirmOrderData{}
	if len(args) == 0 {
		return listResult, nil
	}
	for _, arg := range argsShopServer {
		listResult[arg.(string)] = ConfirmOrderData{
			Jan: jan,
		}
	}
	query := `
SELECT
	wtolf_shop_cd,
	wtolf_server_name,
	wtolf_order_count,
	wtolf_flg_auto_order,
	wtolf_flg_order_export
FROM work_trans_order_list_fix
WHERE TRUE
	AND wtolf_create_date = ?
	AND CONCAT(wtolf_shop_cd,wtolf_server_name,wtolf_franchise_cd,wtolf_jan) IN (` + strings.Trim(strings.Repeat("?,", len(args)), ",") + `)
`
	arg := []interface{}{Common.CurrentDate()}
	for i, _ := range args {
		arg = append(arg, args[i])
	}
	rows, err := this.DB.Query(query, arg...)
	defer rows.Close()
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	for rows.Next() {
		newItem := ConfirmOrderData{}
		err := db.SqlScanStruct(rows, &newItem)
		if err != nil {
			return listResult, exterror.WrapExtError(err)
		}
		listResult[newItem.ShopCd+newItem.ServerName] = newItem
	}

	query = `
SELECT
	moa_server_name wtolf_server_name,
	moa_franchise_cd wtolf_franchise_cd,
	moa_shop_cd wtolf_shop_cd,
	mw_flg_edi wtolf_flg_edi,
	moa_supplier_cd wtolf_order_supplier_cd,
	moa_supplier_line_no wtolf_order_shop_line_no,
	moa_supplier_shop_cd wtolf_order_shop_cd
FROM
	ao_master_order_accounts
JOIN
	ao_master_wholeseller
 ON moa_supplier_cd = mw_supplier_cd
AND moa_franchise_cd = mw_franchise_cd
WHERE TRUE
	AND CONCAT(moa_shop_cd,moa_server_name) IN (` + strings.Trim(strings.Repeat("?,", len(argsShopServer)), ",") + `)
`
	rowsMoa, err := this.DB.Query(query, argsShopServer...)
	defer rowsMoa.Close()
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	for rowsMoa.Next() {
		newItem := ConfirmOrderData{}
		err := db.SqlScanStruct(rowsMoa, &newItem)
		if err != nil {
			return listResult, exterror.WrapExtError(err)
		}
		item := listResult[newItem.ShopCd+newItem.ServerName]
		item.FlgEdi = newItem.FlgEdi
		item.OrderSupplierCd = newItem.OrderSupplierCd
		item.OrderShopLineNo = newItem.OrderShopLineNo
		item.OrderShopCd = newItem.OrderShopCd
		listResult[newItem.ShopCd+newItem.ServerName] = item
	}

	query = `
SELECT
	bqct_servername wtolf_shop_cd,
	bqct_franchise_id wtolf_franchise_cd,
	bqct_shop_cd wtolf_shop_cd,
	genre_cd,
	genre_name,
	bqct_jan_cd
FROM (
SELECT
	bqct_servername,
	bqct_franchise_id,
	bqct_shop_cd,
	bqct_category_cd genre_cd,
	bqct_media_group_name genre_name,
	bqct_jan_cd
FROM bq_category
WHERE bqct_category_type IN ('1','2')
	AND CONCAT(bqct_shop_cd,bqct_servername) IN (` + strings.Trim(strings.Repeat("?,", len(argsShopServer)), ",") + `)
UNION
SELECT
	bqct_servername,
	bqct_franchise_id,
	bqct_shop_cd,
	IF(ssr_refer_genre_type IN ('1','2') , bqct_kubn_cd, bqct_bumon_cd) genre_cd,
	IF(ssr_refer_genre_type IN ('1','2') , bqct_kubn_nm, bqct_bumon_nm) genre_name ,
	bqct_jan_cd
FROM (
SELECT
	bqct_servername,
	bqct_franchise_id,
	bqct_category_type,
	bqct_shop_cd,
	bqct_kubn_cd, bqct_kubn_nm,
	bqct_bumon_cd, bqct_bumon_nm,
	bqct_jan_cd
FROM bq_category
WHERE TRUE
	AND bqct_category_type IN ('3','4')
	AND CONCAT(bqct_shop_cd,bqct_servername) IN (` + strings.Trim(strings.Repeat("?,", len(argsShopServer)), ",") + `)
) bqct
LEFT JOIN setting_shop_refer
	ON bqct_servername = ssr_server_name
	AND bqct_shop_cd = ssr_shop_cd
) bq
WHERE TRUE
	AND CONCAT(bqct_shop_cd,bqct_servername) IN (` + strings.Trim(strings.Repeat("?,", len(argsShopServer)), ",") + `)
	AND bqct_jan_cd = ?
GROUP BY
	bqct_shop_cd,
	bqct_servername,
	bqct_jan_cd
`
	argsBqct := []interface{}{}
	argsBqct = append(argsBqct, argsShopServer...)
	argsBqct = append(argsBqct, argsShopServer...)
	argsBqct = append(argsBqct, argsShopServer...)
	argsBqct = append(argsBqct, jan)

	rowsBqct, err := this.DB.Query(query, argsBqct...)
	defer rowsBqct.Close()
	if err != nil {
		return listResult, exterror.WrapExtError(err)
	}
	for rowsBqct.Next() {
		newItem := ConfirmOrderData{}
		err := db.SqlScanStruct(rowsBqct, &newItem)
		if err != nil {
			return listResult, exterror.WrapExtError(err)
		}
		item := listResult[newItem.ShopCd+newItem.ServerName]
		item.GenreCd = newItem.GenreCd
		item.GenreName = newItem.GoodsName
		listResult[newItem.ShopCd+newItem.ServerName] = item
	}
	return listResult, nil
}

// Insert
// Get receipt in day
func (this *WorkTransOrderListFixModel) MaxReceiptNoInDay() int {
	query := `
SELECT
	IFNULL(MAX(wtolf_receipt_no), 0) max_receipt_no
FROM
	work_trans_order_list_fix
WHERE
	wtolf_create_date = ?
	`
	receiptNo := 0
	rows, err := this.DB.Query(query, Common.CurrentDate())
	defer rows.Close()
	if err != nil {
		return 0
	}
	for rows.Next() {
		rows.Scan(&receiptNo)
	}
	return receiptNo
}

// Insert data
func (this *WorkTransOrderListFixModel) InsertOrderData(args []interface{}) error {
	err := errors.New("インターナルサーバエラー。")
	if len(args) > 0 {
		listField := []string{
			"wtolf_create_date",
			"wtolf_update_date",
			"wtolf_server_name",
			"wtolf_franchise_cd",
			"wtolf_shop_chain_cd",
			"wtolf_shop_cd",
			"wtolf_jan",
			"wtolf_receipt_no",
			"wtolf_order_count",
			"wtolf_order_way_cd",
			"wtolf_flg_edi",
			"wtolf_order_supplier_cd",
			"wtolf_order_shop_line_no",
			"wtolf_order_shop_cd",
			"wtolf_flg_data_fix_date",
			"wtolf_flg_order_export",
			"wtolf_ope_status_fls",
			"wtolf_publisher_cd",
			"wtolf_publisher_name",
			"wtolf_genre_cd",
			"wtolf_genre_name",
			"wtolf_Bgenre_cd",
			"wtolf_location_cd",
		}
		strTemplate := fmt.Sprintf(" ( %v ) ,", Common.SQLPara(listField))
		strValue := strings.Trim(strings.Repeat(strTemplate, len(args)/len(listField)), ",")
		query := `
INSERT INTO
	work_trans_order_list_fix (
	` + strings.Join(listField, ",") + `
	)
VALUE ` + strValue + `
ON DUPLICATE KEY UPDATE
	wtolf_update_date = VALUES(wtolf_update_date),
	wtolf_order_count = VALUES(wtolf_order_count)
`
		_, err = this.DB.Exec(query, args...)
		if err != nil {
			return exterror.WrapExtError(err)
		}
	}
	return nil
}

func (this *WorkTransOrderListFixModel) Search(franchiseCd string, listShopCd []string, orderDateFrom string, orderDateTo string,
	supplierName string, publisherName string, autoFlag string, ediFlag string) chan WtOrderListDataWithErr {

	outputChan := make(chan WtOrderListDataWithErr)

	go func() {
		listQueryArgs := []interface{}{}
		shopQuestions := ""

		listQueryArgs = append(listQueryArgs, franchiseCd)
		for _, v := range listShopCd {
			shopQuestions = shopQuestions + "?,"
			listQueryArgs = append(listQueryArgs, v)
		}
		shopQuestions = strings.Trim(shopQuestions, ",")

		sql := `
SELECT
	LEFT(IFNULL(wtolf_date_fix_date, ''),8) wtolf_flg_data_fix_date,
	IFNULL(mw_supplier_name, '') mw_supplier_name,
	IFNULL(wtolf_order_shop_line_no, '') wtolf_order_shop_line_no,
	IFNULL(wtolf_order_shop_cd, '') wtolf_order_shop_cd,
	IFNULL(wtolf_order_memo, '') wtolf_order_memo,
	IFNULL(wtolf_jan, '') wtolf_jan,
	IFNULL(gm_goods_name, '') gm_goods_name,
	IFNULL(wtolf_publisher_name, '') wtolf_publisher_name,
	IFNULL(gm_author, '') gm_author,
	IFNULL(gm_price, 0) gm_price,
	IFNULL(wtolf_order_count, 0) wtolf_order_count,
	IFNULL(wtolf_shop_cd, '') wtolf_shop_cd,
	IFNULL(shm_shop_name, '') shm_shop_name,
	IFNULL(wtolf_genre_cd, '') wtolf_genre_cd,
	IFNULL(IF(shm_product_cd = 'B', bqct_bumon_nm, IF(shm_product_cd = 'M', bqct_media_group2_name, '')), '') genre_name,
	IF(wtolf_flg_edi = '0', '非データ発注', IF(wtolf_flg_edi = '1', 'データ発注', '')) wtolf_flg_edi,
	IFNULL(mow_order_way_name, '') mow_order_way_name
FROM work_trans_order_list_fix wt
LEFT JOIN ao_master_wholeseller mw
	ON wt.wtolf_order_supplier_cd = mw.mw_supplier_cd
	AND mw_franchise_cd = ?
LEFT JOIN ao_goods_master gm
	ON wt.wtolf_jan = gm.gm_jan
LEFT JOIN bq_category bqct
	ON wt.wtolf_server_name = bqct.bqct_servername
	AND wt.wtolf_shop_cd = bqct.bqct_shop_cd
	AND wt.wtolf_genre_cd = bqct.bqct_category_cd
LEFT JOIN ao_master_order_way mow
	ON wt.wtolf_order_way_cd = mow.mow_order_way_cd
LEFT JOIN shop_master_show shm
	ON wt.wtolf_shop_cd = shm.shm_shop_cd
	AND wt.wtolf_server_name = shm.shm_server_name
WHERE
	CONCAT(wtolf_server_name,'-',wtolf_shop_cd) IN (` + shopQuestions + `)
	AND LEFT(wtolf_date_fix_date, 8) BETWEEN ? AND ?
	`
		listQueryArgs = append(listQueryArgs, orderDateFrom, orderDateTo)

		if len(supplierName) > 0 {
			sql = sql + ` AND IFNULL(mw_supplier_name, '') LIKE ?`
			listQueryArgs = append(listQueryArgs, "%"+supplierName+"%")
		}

		if len(publisherName) > 0 {
			sql = sql + ` AND IFNULL(wtolf_publisher_name, '') LIKE ?`
			listQueryArgs = append(listQueryArgs, "%"+publisherName+"%")
		}

		if autoFlag == flag_manual {
			sql = sql + ` AND wtolf_flg_auto_order = '1'`
		}

		if autoFlag == flag_auto {
			sql = sql + ` AND wtolf_flg_auto_order = '0'`
		}

		if ediFlag == flag_edi_data {
			sql = sql + ` AND wtolf_flg_edi = '1'`
		}

		if ediFlag == flag_edi_nodata {
			sql = sql + ` AND wtolf_flg_edi = '0'`
		}

		rows, err := this.DB.Query(sql, listQueryArgs...)
		if err != nil {
			outputChan <- WtOrderListDataWithErr{nil, exterror.WrapExtError(err)}
			close(outputChan)
			return
		}
		defer rows.Close()

		for rows.Next() {
			newOrderData := WtOrderListData{}
			err := rows.Scan(
				&newOrderData.OrderDate,
				&newOrderData.SupplierName,
				&newOrderData.OrderShopLineNo,
				&newOrderData.OrderShopCd,
				&newOrderData.OrderMemo,
				&newOrderData.ISBN,
				&newOrderData.GoodsName,
				&newOrderData.MakerName,
				&newOrderData.Author,
				&newOrderData.Price,
				&newOrderData.OrderCount,
				&newOrderData.ShopCD,
				&newOrderData.ShopName,
				&newOrderData.GenreCd,
				&newOrderData.GenreName,
				&newOrderData.EdiFlag,
				&newOrderData.OrderWayName,
			)

			if err != nil {
				outputChan <- WtOrderListDataWithErr{nil, exterror.WrapExtError(err)}
				close(outputChan)
				return
			}

			outputChan <- WtOrderListDataWithErr{&newOrderData, nil}
		}

		close(outputChan)
	}()

	return outputChan
}
