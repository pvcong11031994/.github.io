package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"fmt"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strings"
	"time"
)

type TopMonthByYearSalesModel struct {
	DB *sql.DB
}

func (this *TopMonthByYearSalesModel) GetTopMonthByYearSales(listShop, listGenre []string, trnDate string) ([]ModelItems.TopMonthByYearSalesItem, error) {

	var args []interface{}
	listTopMonthByYearSales := []ModelItems.TopMonthByYearSalesItem{}
	condition := ""
	if len(listShop) > 0 && trnDate != "" && len(listGenre) > 0 {
		// Condition transaction date
		condition += `
		 AND tmys_date = ?
		`
		args = append(args, trnDate)

		// Condition list shop
		condition += `
		 AND CONCAT(tmys_servername,'|',tmys_shop_cd) IN(?` + strings.Repeat(",?", len(listShop)-1) + `)
		`
		for _, shopCd := range listShop {
			args = append(args, shopCd)
		}

		// Condition list genre
		condition += `
		 AND tmys_genre_cd IN(?` + strings.Repeat(",?", len(listGenre)-1) + `)
		`
		for _, genreCd := range listGenre {
			args = append(args, genreCd)
		}
	} else {
		return listTopMonthByYearSales, nil
	}
	query := `
SELECT *
FROM (
	SELECT
		tmys_date,
		tmys_year,
		tmys_month,
		tmys_day,
		SUM(tmys_sales_amount) tmys_sales_amount,
		tmys_genre_cd,
		tmys_genre_name
	FROM top_month_by_year_sales
	WHERE TRUE
	` + condition + `
	GROUP BY tmys_genre_cd
	) tmys
WHERE
	tmys_sales_amount > 0
ORDER BY
	tmys_genre_cd
	`
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	// Read value
	for rows.Next() {
		item := ModelItems.TopMonthByYearSalesItem{}
		err = rows.Scan(
			&item.Date,
			&item.Year,
			&item.Month,
			&item.Day,
			&item.SalesAmount,
			&item.GenreCd,
			&item.GenreName,
		)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listTopMonthByYearSales = append(listTopMonthByYearSales, item)
	}

	return listTopMonthByYearSales, nil
}

const _LIMIT = 10

func (this *TopMonthByYearSalesModel) LastDateSalesByDepartment(shop, genre []string) (*[]ModelItems.CompareLastYearByPercentItem, error) {

	condition := ""
	args := []interface{}{}

	if len(shop) > 0 {
		condition += `
		 AND CONCAT(last_date.tmys_servername,'|',last_date.tmys_shop_cd) IN(` + Common.SQLPara(shop) + `)
		`
		args = append(args, Common.ToInterfaceArray(shop)...)
	}

	if len(genre) > 0 {
		condition += `
		 AND last_date.tmys_genre_cd IN(` + Common.SQLPara(genre) + `)
		`
		args = append(args, Common.ToInterfaceArray(genre)...)
	}

	if len(shop) == 0 || len(genre) == 0 {
		return &[]ModelItems.CompareLastYearByPercentItem{}, nil
	}

	query := `
SELECT
	last_date.tmys_genre_cd tmys_genre_cd,
	last_date.tmys_genre_name tmys_genre_name,
	SUM(IFNULL(last_date.tmys_sales_amount,0)) tmys_last_date_sales_amount,
	SUM(IFNULL(last_year.tmys_sales_amount,0)) tmys_last_year_sales_amount
FROM
(
	SELECT *
	FROM top_month_by_year_sales
	WHERE TRUE
	    AND tmys_date = ?
		AND tmys_genre_cd <> ''
) last_date
JOIN
	( SELECT * FROM top_month_by_year_sales WHERE tmys_date = ? ) last_year
	ON  last_date.tmys_servername 	= last_year.tmys_servername
	AND last_date.tmys_dbname 		= last_year.tmys_dbname
	AND last_date.tmys_franchise_cd = last_year.tmys_franchise_cd
	AND last_date.tmys_chain_cd 	= last_year.tmys_chain_cd
	AND last_date.tmys_shop_cd 		= last_year.tmys_shop_cd
	AND last_date.tmys_genre_cd 	= last_year.tmys_genre_cd
WHERE TRUE
` + condition + `
GROUP BY
	last_date.tmys_genre_cd
ORDER BY
	SUM(IFNULL(last_date.tmys_sales_amount,0)) DESC
`

	return this.getData(query, args)
}

func (this *TopMonthByYearSalesModel) LastDateSalesByShop(genre []string) (*[]ModelItems.CompareLastYearByPercentItem, error) {

	condition := ""
	args := []interface{}{}

	if len(genre) > 0 {
		condition += `
		 AND last_date.tmys_genre_cd IN(` + Common.SQLPara(genre) + `)
		`
		args = append(args, Common.ToInterfaceArray(genre)...)
	} else {
		return &[]ModelItems.CompareLastYearByPercentItem{}, nil
	}

	query := `
SELECT
	last_date.tmys_shop_cd tmys_shop_cd,
	shm.shm_shop_name tmys_shop_name,
	SUM(IFNULL(last_date.tmys_sales_amount,0)) tmys_last_date_sales_amount,
	SUM(IFNULL(last_year.tmys_sales_amount,0)) tmys_last_year_sales_amount
FROM
(
	SELECT *
	FROM top_month_by_year_sales
	WHERE TRUE
	    AND tmys_date = ?
		AND tmys_genre_cd <> ''
) last_date
JOIN
	( SELECT * FROM top_month_by_year_sales WHERE tmys_date = ? ) last_year
	ON  last_date.tmys_servername 	= last_year.tmys_servername
	AND last_date.tmys_dbname 		= last_year.tmys_dbname
	AND last_date.tmys_franchise_cd = last_year.tmys_franchise_cd
	AND last_date.tmys_chain_cd 	= last_year.tmys_chain_cd
	AND last_date.tmys_shop_cd 		= last_year.tmys_shop_cd
	AND last_date.tmys_genre_cd 	= last_year.tmys_genre_cd
LEFT JOIN
(
	SELECT shm_server_name, shm_shop_cd, shm_shop_name
	FROM shop_master_show
) shm
	ON  shm.shm_server_name = last_date.tmys_servername
	AND shm.shm_shop_cd     = last_date.tmys_shop_cd
WHERE TRUE
` + condition + `
GROUP BY
	last_date.tmys_servername,
	last_date.tmys_shop_cd
ORDER BY
	SUM(IFNULL(last_date.tmys_sales_amount,0)) DESC
`

	return this.getData(query, args)
}

func (this *TopMonthByYearSalesModel) getData(query string, argsRequest []interface{}) (*[]ModelItems.CompareLastYearByPercentItem, error) {
	query += " LIMIT ? "
	argsRequest = append(argsRequest, _LIMIT)

	// Init date
	lastDate := time.Now().AddDate(0, 0, -1)
	dateOfLastYear := lastDate.AddDate(-1, 0, 0)
	args := []interface{}{
		lastDate.Format(Common.DATE_FORMAT_YMD),
		dateOfLastYear.Format(Common.DATE_FORMAT_YMD),
	}
	args = append(args, argsRequest...)

	listResult := []ModelItems.CompareLastYearByPercentItem{}
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return &listResult, exterror.WrapExtError(err)
	}
	defer rows.Close()

	// Read value
	for rows.Next() {
		item := ModelItems.CompareLastYearByPercentItem{}
		err = gf.SqlScanStruct(rows, &item)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}

		if item.LastYearSalesAmount > 0 {
			percent := float64(item.LastDateSalesAmount) / float64(item.LastYearSalesAmount)
			item.PercentSalesAmount = fmt.Sprintf("%d%v", int(percent*100+0.5), "%")
		} else {
			item.PercentSalesAmount = "～%"
		}
		listResult = append(listResult, item)
	}

	return &listResult, nil
}
