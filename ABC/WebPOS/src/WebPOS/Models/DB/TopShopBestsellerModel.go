package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
	"strings"
)

type TopShopBestsellerModel struct {
	DB *sql.DB
}

func (this *TopShopBestsellerModel) GetTopTenShopBestseller(listShop []string) ([]ModelItems.TopShopBestsellerItem, error) {

	var args []interface{}
	listShopBestseller := []ModelItems.TopShopBestsellerItem{}
	condition := ""
	if len(listShop) > 0 {
		condition += `
		 AND CONCAT(tsb_servername,'|',tsb_shop_cd) IN(?` + strings.Repeat(",?", len(listShop)-1) + `)
		`
		for _, shopCd := range listShop {
			args = append(args, shopCd)
		}
	} else {
		return listShopBestseller, nil
	}
	query := `
SELECT
	tsb_jan,
	tsb_genre_cd,
	tsb_genre_name,
	tsb_goods_name,
	tsb_artist_name,
	tsb_maker_name,
	tsb_price,
	SUM(tsb_sales_count) sales_count,
	SUM(tsb_sales_amount) sales_amount,
	SUM(tsb_stock_count) stock_count
FROM
	top_shop_bestseller
WHERE
	TRUE
	` + condition + `
GROUP BY
	tsb_jan,
	tsb_genre_cd
ORDER BY
	sales_count DESC,
	sales_amount DESC
LIMIT 10
	`
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	// Read value
	for rows.Next() {
		item := ModelItems.TopShopBestsellerItem{}
		err = rows.Scan(
			&item.JanCd,
			&item.GenreCd,
			&item.GenreName,
			&item.GoodsName,
			&item.ArtistName,
			&item.MakerName,
			&item.Price,
			&item.SalesCount,
			&item.SalesAmount,
			&item.Stock,
		)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listShopBestseller = append(listShopBestseller, item)
	}

	return listShopBestseller, nil
}
