package ModelItems

type WidgetMasterItem struct {
	CreateDate string `sql:"wm_create_date"`
	UpdateDate string `sql:"wm_update_date"`
	WidgetId   string `sql:"wm_widget_id"`
	WidgetName string `sql:"wm_widget_name"`
}
