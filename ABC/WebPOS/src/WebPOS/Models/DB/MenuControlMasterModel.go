package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
)

type MenuControlMasterModel struct {
	DB *sql.DB
}

func (this *MenuControlMasterModel) GetMenuName(menu_id, menuGroupFlag string) (ModelItems.MenuLevelItem, error) {

	sqlString := `
	SELECT
		menu_id,
		menu_name
	FROM
		menu_master
	WHERE
		menu_url = ?
	AND menu_group_flg = ?
	`
	item := ModelItems.MenuLevelItem{}
	rows, err := this.DB.Query(sqlString, menu_id, menuGroupFlag)
	if err != nil {
		return item, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&item.MenuId,
			&item.MenuName,
		)
	}
	return item, nil
}

func (this *MenuControlMasterModel) GetMenuLevel(menuGroupFlag string) ([]ModelItems.MenuLevelItem, error) {

	sqlString := `
	SELECT
		menu_id,
		menu_name,
		menu_level,
		parent_menu_id,
		menu_seq_number,
		menu_url
FROM menu_master
WHERE menu_group_flg = ?
ORDER BY
		menu_level,
		menu_seq_number
	`
	item := []ModelItems.MenuLevelItem{}
	rows, err := this.DB.Query(sqlString, menuGroupFlag)
	if err != nil {
		return item, err
	}
	defer rows.Close()
	for rows.Next() {
		items := ModelItems.MenuLevelItem{}
		err = rows.Scan(
			&items.MenuId,
			&items.MenuName,
			&items.MenuLevel,
			&items.ParentMenuId,
			&items.MenuSeqNumber,
			&items.MenuUrl,
		)
		item = append(item, items)
	}
	return item, nil
}

func (this *MenuControlMasterModel) CheckMenuByUrlByFlg(url, menuGroupFlag string) bool {

	query := `
		SELECT   menu_id
				,menu_name
		 FROM menu_master
		 WHERE menu_url = ?
		 AND menu_group_flg = ?
	`
	list_menu := []MenuGroup{}
	rows, err := this.DB.Query(query, url, menuGroupFlag)
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		menu := MenuGroup{}
		err = rows.Scan(
			&menu.Menu_ID,
			&menu.Menu_Name,
		)
		list_menu = append(list_menu, menu)
	}

	if len(list_menu) > 0 {
		return true
	}

	return false
}
