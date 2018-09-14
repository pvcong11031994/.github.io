package ModelItems

type ControlReportDesignItem struct {
	CreateDate        string `sql:"mcrd_create_date"`
	UpdateDate        string `sql:"mcrd_update_date"`
	ShopChainCd       string `sql:"mcrd_shop_chain_cd"`
	ShopCd            string `sql:"mcrd_shop_cd"`
	ColorServiceBar   string `sql:"mcrd_color_service_bar"`
	ColorDashboardBar string `sql:"mcrd_color_dashboard_bar"`
	LogoServiceBar    string `sql:"mcrd_logo_service_bar"`
}
