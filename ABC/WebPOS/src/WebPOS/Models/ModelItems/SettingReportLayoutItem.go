package ModelItems

type SettingReportLayoutItem struct {
	ReportId    string
	ReportName  string
	UserId      string
	SelectedCol []string
	SelectedRow []string
	SelectedSum []string
	MenuId      int
}

type ReportMenuItem struct {
	ReportId   string `sql:"rls_report_id"`
	ReportName string `sql:"rls_report_name"`
	MenuId     string `sql:"rls_report_menu_id"`
}
