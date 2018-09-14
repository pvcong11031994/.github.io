package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"strings"
)

type SettingReportLayoutModel struct {
	DB *sql.DB
}

func (this *SettingReportLayoutModel) GetUserReportMenu(userId string) ([]ModelItems.ReportMenuItem, error) {
	sqlQuery := `
SELECT
	rls_report_id,
	rls_report_name,
	rls_report_menu_id
FROM setting_report_layout
WHERE
	rls_user_id = ?
	AND rls_report_name != ''
ORDER BY rls_report_id, rls_report_menu_id
	`
	rows, err := this.DB.Query(sqlQuery, userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	var listReportMenuItem []ModelItems.ReportMenuItem
	for rows.Next() {
		item := ModelItems.ReportMenuItem{}
		err := db.SqlScanStruct(rows, &item)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}
		listReportMenuItem = append(listReportMenuItem, item)
	}

	return listReportMenuItem, nil
}

func (this *SettingReportLayoutModel) GetSettingReportLayout(userId string, reportId string) (*ModelItems.SettingReportLayoutItem, error) {
	var layout = ModelItems.SettingReportLayoutItem{
		ReportId: reportId,
		UserId:   userId,
	}

	sqlQuery := `
SELECT
	rls_selected_col,
	rls_selected_row,
	rls_selected_sum
FROM setting_report_layout
WHERE
	rls_user_id = ?
	AND rls_report_id = ?
	AND rls_report_name = ''
ORDER BY rls_report_menu_id DESC
LIMIT 1
	`
	rows := this.DB.QueryRow(sqlQuery, userId, reportId)

	var selCols, selRows, selSums string

	err := rows.Scan(&selCols, &selRows, &selSums)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			break
		default:
			return nil, exterror.WrapExtError(err)
		}
	}

	layout.SelectedCol = strings.Split(selCols, ",")
	layout.SelectedRow = strings.Split(selRows, ",")
	layout.SelectedSum = strings.Split(selSums, ",")

	return &layout, nil
}

func (this *SettingReportLayoutModel) GetSettingReportLayoutByMenu(userId string, reportId string, menuId int) (*ModelItems.SettingReportLayoutItem, error) {
	var layout = ModelItems.SettingReportLayoutItem{
		ReportId: reportId,
		UserId:   userId,
	}

	sqlQuery := `
SELECT
	rls_selected_col,
	rls_selected_row,
	rls_selected_sum,
	rls_report_name
FROM setting_report_layout
WHERE
	rls_user_id = ?
	AND rls_report_id = ?
	AND rls_report_menu_id = ?
	AND rls_report_name != ''
ORDER BY rls_report_menu_id DESC
LIMIT 1
	`
	rows := this.DB.QueryRow(sqlQuery, userId, reportId, menuId)

	var selCols, selRows, selSums string

	err := rows.Scan(&selCols, &selRows, &selSums, &layout.ReportName)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			break
		default:
			return nil, exterror.WrapExtError(err)
		}
	}

	layout.SelectedCol = strings.Split(selCols, ",")
	layout.SelectedRow = strings.Split(selRows, ",")
	layout.SelectedSum = strings.Split(selSums, ",")

	return &layout, nil
}

func (this *SettingReportLayoutModel) DeleteReportMenu(userId string, reportId string, menuId int) error {

	sqlQuery := `
DELETE
FROM setting_report_layout
WHERE
	rls_user_id = ?
	AND rls_report_id = ?
	AND rls_report_menu_id = ?
	`
	_, err := this.DB.Exec(sqlQuery, userId, reportId, menuId)

	return exterror.WrapExtError(err)
}

func (this *SettingReportLayoutModel) SetSettingReportLayout(layout *ModelItems.SettingReportLayoutItem) error {

	if layout.ReportName == "" {
		menuId := sql.NullInt64{}
		sqlQuery := `
SELECT
	rls_report_menu_id
FROM setting_report_layout
WHERE
	rls_user_id = ?
	AND rls_report_id = ?
	AND rls_report_name = ''
ORDER BY rls_report_menu_id DESC
LIMIT 1
	`
		row := this.DB.QueryRow(sqlQuery, layout.UserId, layout.ReportId)
		row.Scan(&menuId)

		sql := `
INSERT INTO setting_report_layout (
	rls_user_id,
	rls_report_id,
	rls_report_menu_id,
	rls_selected_col,
	rls_selected_row,
	rls_selected_sum)
VALUES(?,?,?,?,?,?)
ON DUPLICATE KEY UPDATE
	rls_selected_col = VALUES(rls_selected_col),
	rls_selected_row = VALUES(rls_selected_row),
	rls_selected_sum = VALUES(rls_selected_sum)
	`
		_, err := this.DB.Exec(sql, layout.UserId, layout.ReportId, menuId,
			strings.Join(layout.SelectedCol, ","),
			strings.Join(layout.SelectedRow, ","),
			strings.Join(layout.SelectedSum, ","))

		return exterror.WrapExtError(err)
	} else {
		sql := `
INSERT INTO setting_report_layout (
	rls_user_id,
	rls_report_id,
	rls_report_name,
	rls_selected_col,
	rls_selected_row,
	rls_selected_sum)
VALUES(?,?,?,?,?,?)
`
		result, err := this.DB.Exec(sql, layout.UserId, layout.ReportId, layout.ReportName,
			strings.Join(layout.SelectedCol, ","),
			strings.Join(layout.SelectedRow, ","),
			strings.Join(layout.SelectedSum, ","))
		if err != nil {
			layout.MenuId = 0
		} else {
			lastInsertId, err := result.LastInsertId()
			if err == nil {
				layout.MenuId = int(lastInsertId)
			}
		}
		return exterror.WrapExtError(err)
	}
}
