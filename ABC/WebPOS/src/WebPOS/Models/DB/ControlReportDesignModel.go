package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type ReportDesignMaster struct {
	DB *sql.DB
}

func (this *ReportDesignMaster) GetReportDesign(chain string) (*ModelItems.ControlReportDesignItem, error) {
	reportDesign := ModelItems.ControlReportDesignItem{}

	sql := `
SELECT
	mcrd_color_service_bar,
	mcrd_color_dashboard_bar,
	mcrd_logo_service_bar
FROM ao_master_control_rep_design
WHERE mcrd_shop_chain_cd = ?
	`

	rows, err := this.DB.Query(sql, chain)
	Common.LogErr(exterror.WrapExtError(err))
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			err = rows.Scan(
				&reportDesign.ColorServiceBar,
				&reportDesign.ColorDashboardBar,
				&reportDesign.LogoServiceBar,
			)
			if err != nil {
				return nil, exterror.WrapExtError(err)
			}
		}
	}
	return &reportDesign, nil
}
