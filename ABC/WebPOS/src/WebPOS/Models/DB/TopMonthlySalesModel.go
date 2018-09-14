package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
	"strings"
)

type TopMonthlySalesModel struct {
	DB *sql.DB
}

func (this *TopMonthlySalesModel) GetTopMonthlySalesByMonth(listShop []string, dateTimeYMD string) ([]ModelItems.TopMonthlySalesItem, error) {

	var args []interface{}

	listTopMonthlySales := []ModelItems.TopMonthlySalesItem{}
	condition := ""
	if len(listShop) > 0 && len(dateTimeYMD) >= 8 {
		// Add param this month and this year
		args = append(args, dateTimeYMD[:4])
		args = append(args, dateTimeYMD[4:6])

		// Condition shop
		condition += `
		 AND CONCAT(tms_servername,'|',tms_shop_cd) IN(?` + strings.Repeat(",?", len(listShop)-1) + `)
		`
		for _, shopCd := range listShop {
			args = append(args, shopCd)
		}
	} else {
		return listTopMonthlySales, nil
	}
	query := `
SELECT
	tms_date,
	tms_year,
	tms_month,
	tms_day,
	SUM(tms_sales_amount) tms_sales_amount
FROM top_monthly_sales
WHERE
	tms_year = ?
	AND tms_month = ?
` + condition + `
GROUP BY tms_date
ORDER BY tms_date
	`

	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	// Read value
	for rows.Next() {
		item := ModelItems.TopMonthlySalesItem{}
		err = rows.Scan(
			&item.TrnDate,
			&item.TrnYear,
			&item.TrnMonth,
			&item.TrnDay,
			&item.SalesAmount,
		)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listTopMonthlySales = append(listTopMonthlySales, item)
	}

	return listTopMonthlySales, nil
}
