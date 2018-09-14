package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
)

type WidgetMasterModel struct {
	DB *sql.DB
}

func (this *WidgetMasterModel) WidgetGroup() (*[]ModelItems.WidgetMasterItem, error) {

	query := `
SELECT
	wm_widget_id,
	wm_widget_name
FROM
	widget_master
ORDER BY
	wm_widget_id
`
	// Execute
	rows, err := this.DB.Query(query)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	//
	listWidgetSetting := []ModelItems.WidgetMasterItem{}
	for rows.Next() {
		item := ModelItems.WidgetMasterItem{}
		if err = gf.SqlScanStruct(rows, &item); err != nil {
			return &listWidgetSetting, exterror.WrapExtError(err)
		}
		listWidgetSetting = append(listWidgetSetting, item)
	}
	return &listWidgetSetting, nil
}
