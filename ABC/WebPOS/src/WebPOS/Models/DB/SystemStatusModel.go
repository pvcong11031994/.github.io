package Models

import (
	"WebPOS/Common"
	"database/sql"
	"github.com/goframework/gf/exterror"
	"WebPOS/Models/ModelItems"
	"time"
)

type SystemStatusModel struct {
	DB *sql.DB
}

func (this *SystemStatusModel) GetStatus(limit int) []ModelItems.SystemStatusItem {
	listStatus := []ModelItems.SystemStatusItem{}

	sql := `
SELECT
	IFNULL(ss_created_at, ''),
	IFNULL(ss_chain, ''),
	IFNULL(ss_group, ''),
	IFNULL(ss_detail, '')
FROM ao_system_status
ORDER BY ss_created_at DESC
LIMIT ?
	`

	rows, err := this.DB.Query(sql, limit)
	Common.LogErr(exterror.WrapExtError(err))
	if err == nil {
		newStatus := ModelItems.SystemStatusItem{}
		for rows.Next() {
			err := rows.Scan(
				&newStatus.CreatedAt,
				&newStatus.Chain,
				&newStatus.Group,
				&newStatus.Detail)
			Common.LogErr(exterror.WrapExtError(err))

			if err == nil {
				listStatus = append(listStatus, newStatus)
			}
		}
		defer rows.Close()
	}

	return listStatus
}

func (this *SystemStatusModel) InsertStatus(chain, group, detail string) error {

	query := `
INSERT INTO ao_system_status (
	ss_created_at,
	ss_chain,
	ss_group,
	ss_detail
) VALUE (?, ?, ?, ?)
	`

	_, err := this.DB.Exec(query, time.Now().Format(Common.DATE_FORMAT_MYSQL_YMDHMS), chain, group, detail)
	if err != nil {
		return exterror.WrapExtError(err)
	}

	return nil
}