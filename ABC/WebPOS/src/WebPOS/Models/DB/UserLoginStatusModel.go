package Models

import (
	"database/sql"
	"github.com/goframework/gf/exterror"
)

type UserLoginStatusModel struct {
	DB *sql.DB
}

// Check status login by user ID and session ID
func (this *UserLoginStatusModel) CheckRejected(userID string, sessionID string) (bool, error) {

	sql := `
SELECT
	COUNT(*)
FROM user_login_status
WHERE
	uls_user_id = ?
	AND uls_login_session_id = ?
	AND uls_login_status = '1'
	`

	rows, err := this.DB.Query(sql, userID, sessionID)
	if err != nil {
		return true, exterror.WrapExtError(err)
	}
	defer rows.Close()
	if rows.Next() {
		count := -1
		rows.Scan(&count)
		if count == 0 {
			return false, nil
		}
	}
	return true, exterror.WrapExtError(err)
}

// Update status when user logout
func (this *UserLoginStatusModel) LogoutByKey(userID string, sessionID string) error {

	sql := `
UPDATE user_login_status
SET
	uls_logout_time = now(),
	uls_logout_reason = '0',
	uls_login_status = '1'
WHERE
	uls_user_id = ?
	AND uls_login_session_id = ?
	AND uls_login_status = '0'
	`
	_, err := this.DB.Exec(sql, userID, sessionID)

	return exterror.WrapExtError(err)
}

// Update status when session timeout or problem of session_store
func (this *UserLoginStatusModel) RejectLoggedInUser(userID string) error {

	sql := `
UPDATE user_login_status
SET
	uls_logout_time = now(),
	uls_logout_reason = '1',
	uls_login_status = '1'
WHERE
	uls_user_id = ?
	AND uls_login_status = '0'
	`
	_, err := this.DB.Exec(sql, userID)

	return exterror.WrapExtError(err)
}

// Check uls_login_status of userId
func (this *UserLoginStatusModel) CheckLoginUser(userID string) (bool, error) {

	sql := `
SELECT
	IF(COUNT(*) > 0, true, false)
FROM
	user_login_status
WHERE
	uls_user_id = ?
	AND uls_login_status = '0'
	`
	isLogin := false
	row := this.DB.QueryRow(sql, userID)
	err := row.Scan(&isLogin)

	return isLogin, exterror.WrapExtError(err)
}

// Add record user_login_status when user login
func (this *UserLoginStatusModel) NewLoginStatus(userID, sessionID, ipAddress, browserInfo, deviceID string) error {

	sql := `
INSERT INTO user_login_status (
	uls_create_date,
	uls_update_date,
	uls_user_id,
	uls_login_session_id,
	uls_login_status,
	uls_login_time,
	uls_login_ip,
	uls_login_browser,
	uls_logout_time,
	uls_logout_reason,
	uls_device_id
) VALUES (
	now(),
	now(),
	?,      -- uls_user_id
	?,      -- uls_login_session_id
	'0',    -- uls_login_status
	now(),  -- uls_login_time
	?,      -- uls_login_ip
	?,      -- uls_login_browser
	NULL,   -- uls_logout_time
	'',     -- uls_logout_reason
	?       -- uls_device_id
)
	`

	_, err := this.DB.Exec(sql, userID, sessionID, ipAddress, browserInfo, deviceID)

	return exterror.WrapExtError(err)
}

// Get session ID of user is login
func (this *UserLoginStatusModel) GetLoginInfoByUser(userID string) (string, error) {

	sql := `
SELECT
	uls_login_session_id
FROM user_login_status
WHERE
	uls_user_id = ?
	AND uls_login_status = '0'
	`

	ssId := ""
	row := this.DB.QueryRow(sql, userID)
	err := row.Scan(&ssId)

	return ssId, exterror.WrapExtError(err)
}
