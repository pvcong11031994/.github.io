package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"time"
)

type TopDailyShopSalesStockModel struct {
	DB *sql.DB
}

func (this *TopDailyShopSalesStockModel) StockSalesInMonth(franchise, shop []string) (*[]ModelItems.TopDailyShopSalesStockItem, error) {

	listResult := []ModelItems.TopDailyShopSalesStockItem{}
	condition := ""
	args := []interface{}{time.Now().AddDate(0, 0, -1).Format("200601")}

	if len(franchise) > 0 {
		condition += `
		 AND tdsss_franchise_cd IN(` + Common.SQLPara(franchise) + `)
		`
		args = append(args, Common.ToInterfaceArray(franchise)...)
	}

	if len(shop) > 0 {
		condition += `
		 AND CONCAT(tdsss_servername,'|',tdsss_shop_cd)  IN(` + Common.SQLPara(shop) + `)
		`
		args = append(args, Common.ToInterfaceArray(shop)...)
	}

	if len(franchise) == 0 || len(shop) == 0 {
		return &[]ModelItems.TopDailyShopSalesStockItem{}, nil
	}

	query := `
SELECT
	tdsss_date,
	SUM(tdsss_goods_count) tdsss_goods_count,
	SUM(tdsss_stock_count) tdsss_stock_count
FROM top_daily_shop_sales_stock
WHERE TRUE
AND LEFT(tdsss_date, 6) = ?
` + condition + `
GROUP BY tdsss_date
ORDER BY tdsss_date
`
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return &listResult, exterror.WrapExtError(err)
	}
	defer rows.Close()

	// Read value
	for rows.Next() {
		item := ModelItems.TopDailyShopSalesStockItem{}
		err = gf.SqlScanStruct(rows, &item)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		date, _ := time.Parse(Common.DATE_FORMAT_YMD, item.Date)
		item.Date = date.Format(Common.JP_DATE_H_MONTH_DAY)
		listResult = append(listResult, item)
	}
	return &listResult, nil
}
