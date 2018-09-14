package Models

import (
	"WebPOS/Common"
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type UserNoticeModel struct {
	DB *sql.DB
}

func (this *UserNoticeModel) GetNotice(chain_cd string) (string, error) {

	query := `
SELECT IFNULL(un_content, "") un_content
FROM ao_user_notice
WHERE un_chain_cd = ?
UNION
SELECT "" un_content
		`

	var content = ""
	err := this.DB.QueryRow(query, chain_cd).Scan(&content)

	return content, exterror.WrapExtError(err)
}

func (this *UserNoticeModel) SetNotice(chain_cd string, user string, content string) error {
	query := `
INSERT INTO ao_user_notice (
	un_created_at,
	un_updated_at,
	un_created_by,
	un_updated_by,
	un_chain_cd,
	un_content)
VALUES (?,?,?,?,?,?)
ON DUPLICATE KEY UPDATE
	un_updated_at = VALUES(un_updated_at),
	un_updated_by = VALUES(un_updated_by),
	un_content = VALUES(un_content)
		`
	now := Common.CurrentDateTime()
	_, err := this.DB.Exec(query, now, now, user, user, chain_cd, content)

	return exterror.WrapExtError(err)
}
