package Models

import (
	"WebPOS/Models/ModelItems"
	"database/sql"
)

type MasterCalendarModel struct {
	DB *sql.DB
}

func (this *MasterCalendarModel) GetDay(dateFrom, dateTo string) ([]ModelItems.MasterCalendarItem, error) {

	sqlString := `
	SELECT
		mc_yyyy,
		mc_mm,
		mc_dd
	FROM
		master_calendar
	WHERE
		mc_yyyymmdd >= ?
		AND mc_yyyymmdd <= ?
	ORDER BY
		mc_yyyy,
		mc_mm,
		mc_dd
	`
	listResult := []ModelItems.MasterCalendarItem{}
	rows, err := this.DB.Query(sqlString, dateFrom, dateTo)
	if err != nil {
		return listResult, err
	}
	defer rows.Close()
	for rows.Next() {
		item := ModelItems.MasterCalendarItem{}
		err = rows.Scan(
			&item.Mcyyyy,
			&item.Mcmm,
			&item.Mcdd,
		)
		item.McKey = item.Mcyyyy + item.Mcmm + item.Mcdd
		listResult = append(listResult, item)
	}
	return listResult, nil
}

func (this *MasterCalendarModel) GetWeek(dateFrom, dateTo string) ([]ModelItems.MasterCalendarItem, error) {

	sqlString := `
	SELECT
		mc_yyyy,
		MAX(mc_mm) mc_mm,
		MAX(mc_dd) mc_dd,
		mc_weeknum,
		MAX(mc_weekdate) mc_weekdate
	FROM
		master_calendar
	WHERE
		mc_yyyymmdd >= ?
		AND mc_yyyymmdd <= ?
	GROUP BY
		mc_yyyy,
		mc_weeknum
	ORDER BY
		mc_yyyy,
		mc_mm,
		mc_dd
	`
	listResult := []ModelItems.MasterCalendarItem{}
	rows, err := this.DB.Query(sqlString, dateFrom, dateTo)
	if err != nil {
		return listResult, err
	}
	defer rows.Close()
	for rows.Next() {
		item := ModelItems.MasterCalendarItem{}
		err = rows.Scan(
			&item.Mcyyyy,
			&item.Mcmm,
			&item.Mcdd,
			&item.Mcweeknum,
			&item.Mcweekdate,
		)
		item.McKey = item.Mcyyyy + item.Mcweeknum
		listResult = append(listResult, item)
	}
	return listResult, nil
}

func (this *MasterCalendarModel) GetMonth(dateFrom, dateTo string) ([]ModelItems.MasterCalendarItem, error) {

	sqlString := `
	SELECT
		mc_yyyy,
		mc_mm
	FROM
		master_calendar
	WHERE
		mc_yyyymmdd >= ?
		AND mc_yyyymmdd <= ?
	GROUP BY
		mc_yyyy, mc_mm
	ORDER BY
		mc_yyyy,
		mc_mm
	`
	listResult := []ModelItems.MasterCalendarItem{}
	rows, err := this.DB.Query(sqlString, dateFrom, dateTo)
	if err != nil {
		return listResult, err
	}
	defer rows.Close()
	for rows.Next() {
		item := ModelItems.MasterCalendarItem{}
		err = rows.Scan(
			&item.Mcyyyy,
			&item.Mcmm,
		)
		item.McKey = item.Mcyyyy + item.Mcmm
		listResult = append(listResult, item)
	}
	return listResult, nil
}