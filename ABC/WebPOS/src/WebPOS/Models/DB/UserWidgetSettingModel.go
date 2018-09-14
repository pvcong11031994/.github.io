package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type UserWidgetSettingModel struct {
	DB *sql.DB
}

func (this *UserWidgetSettingModel) GetUserWidgetSettingByUserId(userId string) ([]ModelItems.UserWidgetSettingItem, error) {

	query := `
SELECT
	uws_user_id,
	uws_location_frame_id,
	uws_widget_id,
	wm_widget_name
FROM user_widget_setting
JOIN widget_master
	ON uws_widget_id = wm_widget_id
WHERE uws_user_id = ?
AND uws_location_frame_id <> ''
	`
	rows, err := this.DB.Query(query, userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	// Read value
	listWidgetSetting := []ModelItems.UserWidgetSettingItem{}
	for rows.Next() {
		item := ModelItems.UserWidgetSettingItem{}
		err = rows.Scan(
			&item.UserId,
			&item.LocationFrameId,
			&item.WidgetId,
			&item.WidgetName,
		)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listWidgetSetting = append(listWidgetSetting, item)
	}

	return listWidgetSetting, nil
}

func (this *UserWidgetSettingModel) InsertWidgetSettingByUserId(userId, widgetId, frameId string) error {
	args := []interface{}{}
	query := ""
	if widgetId == "" {
		query = `
DELETE
FROM user_widget_setting
WHERE TRUE
AND uws_user_id = ?
AND uws_location_frame_id = ?
`
		args = append(args, userId, frameId)
	} else {
		timeNow := Common.CurrentDate()
		query = `
INSERT INTO user_widget_setting (
	uws_create_date,
	uws_update_date,
	uws_user_id,
	uws_widget_id,
	uws_location_frame_id)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	uws_widget_id = VALUES(uws_widget_id),
	uws_update_date = VALUES(uws_update_date)
`
		args = append(args, timeNow, timeNow, userId, widgetId, frameId)
	}
	_, err := this.DB.Exec(query, args...)
	return exterror.WrapExtError(err)
}
