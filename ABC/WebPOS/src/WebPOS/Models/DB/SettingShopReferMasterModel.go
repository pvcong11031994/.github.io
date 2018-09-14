package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type SettingShopReferMasterModel struct {
	DB *sql.DB
}

// Using in check 参照部門テーブル
func (this *SettingShopReferMasterModel) GetInfoGenreReferByShop(shop string, referMaster *ModelItems.SettingShopReferItem) error {
	sqlString := `
SELECT
	ssr_create_date,
	ssr_update_date,
	ssr_shop_cd,
	ssr_refer_genre_master,
	ssr_refer_genre_type,
	ssr_start_genre_cd,
	ssr_end_genre_cd,
	ssr_genre_cd_check_off,
	ssr_toritsugi_cd1,
	ssr_toritsugi_cd2,
	ssr_toritsugi_cd3,
	ssr_priority_gm
FROM
	setting_shop_refer
WHERE
	CONCAT(ssr_server_name,'|',ssr_shop_cd) = ?
	`

	err := this.DB.QueryRow(sqlString, shop).Scan(
		&referMaster.CreateDate,
		&referMaster.UpdateDate,
		&referMaster.ShopCd,
		&referMaster.ReferGenreMaster,
		&referMaster.ReferGenreType,
		&referMaster.StartGenreCd,
		&referMaster.EndGenreCd,
		&referMaster.GenreCdCheckOff,
		&referMaster.ToritsugiCd1,
		&referMaster.ToritsugiCd2,
		&referMaster.ToritsugiCd3,
		&referMaster.PriorityGm,
	)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	return nil
}
