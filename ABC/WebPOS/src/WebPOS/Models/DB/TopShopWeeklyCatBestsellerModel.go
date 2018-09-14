package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"time"
)

type TopShopWeeklyCatBestsellerItem struct {
	DB *sql.DB
}

const _LIMIT_TOP_JAN_WEEK = 10

func (this *TopShopWeeklyCatBestsellerItem) WeeklyBestbyDepartment(franchise, shop, genre []string) (*map[string]map[string]interface{}, *[]ModelItems.TopShopWeeklyCatBestsellerItem, error) {

	condition := ""
	args := []interface{}{_LIMIT_TOP_JAN_WEEK}

	if len(franchise) > 0 {
		condition += `
		 AND tswcb_franchise_cd IN(` + Common.SQLPara(franchise) + `)
		`
		args = append(args, Common.ToInterfaceArray(franchise)...)
	}
	if len(shop) > 0 {
		condition += `
		 AND CONCAT(tswcb_servername,'|',tswcb_shop_cd)  IN(` + Common.SQLPara(shop) + `)
		`
		args = append(args, Common.ToInterfaceArray(shop)...)
	}
	if len(genre) > 0 {
		condition += `
		 AND tswcb_genre_cd IN(` + Common.SQLPara(genre) + `)
		`
		args = append(args, Common.ToInterfaceArray(genre)...)
	}

	query := `
SELECT
	_data.tswcb_jan,
	_date.tswcb_date,
	MAX(_data.tswcb_artist_name) tswcb_artist_name,
	MAX(_data.tswcb_maker_name) tswcb_maker_name,
	SUM(IFNULL(_data.tswcb_sales_count,0)) tswcb_sales_count
FROM
(
SELECT DISTINCT tswcb_date
FROM top_shop_weekly_cat_bestseller
) _date
LEFT JOIN (
	SELECT *
	FROM top_shop_weekly_cat_bestseller
	WHERE tswcb_jan IN (
		SELECT 	tswcb_jan
		FROM (
			SELECT 	tswcb_jan,
					SUM(IFNULL(tswcb_sales_count,0)) tswcb_sales_count
			FROM top_shop_weekly_cat_bestseller
			GROUP BY tswcb_jan
			ORDER BY
				SUM(IFNULL(tswcb_sales_count,0)) DESC
		) limit_jan
	)
` + condition + `
LIMIT ?
) _data
ON _date.tswcb_date = _data.tswcb_date
GROUP BY
	_date.tswcb_date,
	_data.tswcb_jan
`
	listResult := []ModelItems.TopShopWeeklyCatBestsellerItem{}
	totalSumData := map[string]map[string]interface{}{}
	if len(franchise) == 0 || len(shop) == 0 || len(genre) == 0 {
		return &totalSumData, &listResult, nil
	}
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return &totalSumData, &listResult, exterror.WrapExtError(err)
	}
	defer rows.Close()

	sumData := map[string]interface{}{}
	distinctJan := map[string]interface{}{}
	distinctDate := map[string]interface{}{}
	countDateJan := map[string]interface{}{}

	// Read value
	for rows.Next() {
		item := ModelItems.TopShopWeeklyCatBestsellerItem{}
		err = gf.SqlScanStruct(rows, &item)
		if err != nil {
			return &totalSumData, &listResult, exterror.WrapExtError(err)
		}
		date, _ := time.Parse(Common.DATE_FORMAT_YMD, item.Date)
		item.Date = date.Format(Common.DATE_FORMAT_MD)

		// Calculator total sales count  by date, jan
		totalDate := 0
		if sumData[item.Date] != nil {
			totalDate = sumData[item.Date].(int)
		}
		totalJan := 0
		if sumData[item.Jan] != nil {
			totalJan = sumData[item.Jan].(int)
		}
		totalEnd := 0
		if sumData["total"] != nil {
			totalEnd = sumData["total"].(int)
		}

		sumData[item.Date] = interface{}(totalDate + item.SalesCount)
		sumData[item.Jan] = interface{}(totalJan + item.SalesCount)
		sumData["total"] = interface{}(totalEnd + item.SalesCount)

		// Remove duplia
		distinctJan[item.Jan] = item.Jan
		distinctDate[item.Date] = item.Date
		countDateJan[item.Date+"_"+item.Jan] = item.SalesCount

		listResult = append(listResult, item)
	}

	totalSumData["sum_data"] = sumData
	totalSumData["date"] = distinctDate
	totalSumData["jan"] = distinctJan
	totalSumData["count_date_jan"] = countDateJan

	return &totalSumData, &listResult, nil
}