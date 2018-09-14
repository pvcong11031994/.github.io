package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
	"strings"
)

type TopDayGenreSalesModel struct {
	DB *sql.DB
}

func (this *TopDayGenreSalesModel) GetTopDayGenreSalesByTrnDate(listShop, listGenre []string, trnDate string) ([]ModelItems.TopDayGenreSalesItem, error) {

	var args []interface{}
	listTopDayGenreSales := []ModelItems.TopDayGenreSalesItem{}
	condition := ""
	if len(listShop) > 0 && trnDate != "" && len(listGenre) > 0 {
		// Condition transaction date
		condition += `
		 AND tpgs_date = ?
		`
		args = append(args, trnDate)

		// Condition list shop
		condition += `
		 AND CONCAT(tpgs_servername,'|',tpgs_shop_cd) IN(?` + strings.Repeat(",?", len(listShop)-1) + `)
		`
		for _, shopCd := range listShop {
			args = append(args, shopCd)
		}

		// Condition list genre
		condition += `
		 AND tpgs_genre_cd IN(?` + strings.Repeat(",?", len(listGenre)-1) + `)
		`
		for _, genreCd := range listGenre {
			args = append(args, genreCd)
		}
	} else {
		return listTopDayGenreSales, nil
	}
	query := `
SELECT *
FROM (
	SELECT
		tpgs_date,
		tpgs_year,
		tpgs_month,
		tpgs_day,
		SUM(tpgs_sales_amount) tpgs_sales_amount,
		tpgs_genre_cd,
		tpgs_genre_name
	FROM top_day_genre_sales
	WHERE TRUE
	` + condition + `
	GROUP BY tpgs_genre_cd
	) tpgs
WHERE
	tpgs_sales_amount > 0
ORDER BY
	tpgs_genre_cd
	`

	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	// Read value
	for rows.Next() {
		item := ModelItems.TopDayGenreSalesItem{}
		err = rows.Scan(
			&item.TrnDate,
			&item.TrnYear,
			&item.TrnMonth,
			&item.TrnDay,
			&item.SalesAmount,
			&item.GenreCd,
			&item.GenreName,
		)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listTopDayGenreSales = append(listTopDayGenreSales, item)
	}

	return listTopDayGenreSales, nil
}
