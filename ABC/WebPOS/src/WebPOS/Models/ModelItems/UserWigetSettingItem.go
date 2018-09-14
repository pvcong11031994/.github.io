package ModelItems

type UserWidgetSettingItem struct {
	UserId          string `sql:"uws_user_id"`
	LocationFrameId string `sql:"uws_location_frame_id"`
	WidgetId        string `sql:"uws_widget_id"`
	WidgetName      string `sql:"wm_widget_name"`
}
