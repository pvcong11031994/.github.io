package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/goframework/gf"
	"github.com/goframework/gf/exterror"
	"strings"
)

const (
	_PASSWORD_LENGTH_DEFAULT = 4
)

var umTBLName = map[string]string{
	_SERVER_DEFAULT: "ao_user_master",
}

type UserMasterAPI_Model struct {
	db        *sql.DB
	tableName map[string]string
}

// Create UserMasterAPI_Model with inject db connect
func NewUserMaster(db *sql.DB) *UserMasterAPI_Model {
	//
	newModel := UserMasterAPI_Model{}
	//
	newModel.db = db
	newModel.tableName = map[string]string{}
	newModel.tableName[":um_tb_name"] = umTBLName[_SERVER_DEFAULT]
	return &newModel
}

/*-------------------------------*/
/*      User master helpful      */
/*-------------------------------*/
// Replace table name
func (this *UserMasterAPI_Model) renderTbName(query string) string {
	//
	for key, name := range this.tableName {
		query = strings.Replace(query, key, name, -1)
	}
	return query
}

// Execute query
func (this *UserMasterAPI_Model) execute(query string, args []interface{}) error {
	//
	query = this.renderTbName(query)
	_, err := this.db.Exec(query, args...)
	//
	return err
}

// Execute query with transaction
func (this *UserMasterAPI_Model) executeTrans(query string, args []interface{}) error {
	//
	query = this.renderTbName(query)
	trans, err := this.db.Begin()
	if err != nil {
		return err
	}
	_, err = trans.Exec(query, args...)
	if err != nil {
		trans.Rollback()
		return err
	}
	//
	trans.Commit()
	return err
}

// Query multi row data
func (this *UserMasterAPI_Model) query(query string, args []interface{}) (*sql.Rows, error) {
	//
	query = this.renderTbName(query)

	return this.db.Query(query, args...)
}

// Query one row data
func (this *UserMasterAPI_Model) queryRow(query string, args []interface{}) *sql.Row {
	//
	query = this.renderTbName(query)
	return this.db.QueryRow(query, args...)
}

// Check exit user
func (this *UserMasterAPI_Model) IsExist(user_id string) (bool, error) {
	//
	queryStr := `
SELECT IF(COUNT(*), true, false)
FROM :um_tb_name
WHERE TRUE
	AND um_user_id = ?
`
	var isExists bool
	err := this.queryRow(queryStr, []interface{}{user_id}).Scan(&isExists)
	//
	return isExists, exterror.WrapExtError(err)
}

/*-----------------------------------------------------------------------------------------------------------------------------*/

// Search
func (this *UserMasterAPI_Model) GetUsers(parameter map[string]interface{}) (*[]ModelItems.UserItem, int, error) {
	//
	args := []interface{}{}
	conditionStr := ""
	// ユーザID
	if parameter["um_user_id"].(string) != "" {
		conditionStr += " AND um_user_id = ? "
		args = append(args, parameter["um_user_id"].(string))
	}
	// ユーザ名
	if parameter["um_user_name"].(string) != "" {
		conditionStr += " AND um_user_name LIKE ? "
		args = append(args, fullLike(parameter["um_user_name"].(string)))
	}
	// 部署コード
	if parameter["um_dept_cd"].(string) != "" {
		conditionStr += " AND um_dept_cd LIKE ? "
		args = append(args, fullLike(parameter["um_dept_cd"].(string)))
	}
	// 稼働フラグ
	if parameter["um_flg_use"].(string) != "" {
		conditionStr += " AND um_flg_use = ? "
		args = append(args, parameter["um_flg_use"].(string))
	}
	// 20170623 ASO-4808#comment-3214533
	// conditionStr += " LIMIT ? OFFSET ? "
	// args = append(args, parameter["LIMIT"].(int), parameter["OFFSET"].(int))
	// 20170623 ASO-4808#comment-3214533 end

	queryStr := `
SELECT *
FROM :um_tb_name
WHERE TRUE
` + conditionStr

	//
	rows, err := this.query(queryStr, args)
	if err != nil {
		return nil, 0, exterror.WrapExtError(err)
	}
	defer rows.Close()
	//
	resultList := []ModelItems.UserItem{}
	// 20170623 ASO-4808#comment-3214533
	limit := parameter["LIMIT"].(int)
	offset := parameter["OFFSET"].(int)
	i := 0
	for rows.Next() {
		if (i+1) > offset && (i+1) <= (offset+limit) {
			newItem := ModelItems.UserItem{}
			err := gf.SqlScanStruct(rows, &newItem)
			if err != nil {
				return nil, 0, exterror.WrapExtError(err)
			}
			resultList = append(resultList, newItem)
		}
		i++
	}
	return &resultList, i, nil
}

// Register
func (this *UserMasterAPI_Model) AddUser(parameter []byte) (string, error) {
	//
	newUser := ModelItems.UserItem{}
	err := json.Unmarshal(parameter, &newUser)
	if err != nil {
		return "", exterror.WrapExtError(err)
	}
	//
	currentDateTime := Common.CurrentDatetime()
	expiringDate := Common.DatetimePastWithFormat(1, Common.DATE_FORMAT_ZERO_TIME)
	randomPass := Common.RandNumber(_PASSWORD_LENGTH_DEFAULT)
	encryptPass := Common.GeneratePass(newUser.UserID, randomPass)
	//
	newUser.CreateDate = currentDateTime
	newUser.UpdateDate = currentDateTime
	newUser.UserPass = encryptPass
	newUser.PwExpiringDate = expiringDate
	newUser.PwChangeFlg = "0"
	if newUser.OptionFlgReturn == "" {
		newUser.OptionFlgReturn = "0"
	}

	//
	args := []interface{}{
		newUser.CreateDate,
		newUser.UpdateDate,
		newUser.UserID,
		newUser.UserName,
		newUser.FlgAuth,
		newUser.ShopChainCd,
		newUser.FranchiseCd,
		newUser.FranchiseGroupCd,
		newUser.ServerName,
		newUser.ShopCd,
		newUser.ShopName,
		newUser.FlgMenuGroup,
		newUser.DeptCd,
		newUser.DeptName,
		newUser.UserMail,
		newUser.UserPhone,
		newUser.UserXerox,
		newUser.UserPass,
		newUser.FlgUse,
		newUser.CorpCd,
		newUser.CorpName,
		newUser.ReferOrderTable,
		newUser.PublisherFlg,
		newUser.PublisherId,
		newUser.PwExpiringDate,
		newUser.PwChangeFlg,
		newUser.OptionFlgReturn,
	}
	strValue := strings.Trim(strings.Repeat("?,", 27), ",")
	strValue = fmt.Sprintf(" ( %v ); ", strValue)
	//
	queryStr := `
INSERT INTO :um_tb_name (
	um_create_date,
	um_update_date,
	um_user_ID,
	um_user_name,
	um_flg_auth,
	um_shop_chain_cd,
	um_franchise_cd,
	um_franchise_group_cd,
	um_server_name,
	um_shop_cd,
	um_shop_name,
	um_flg_menu_group,
	um_dept_cd,
	um_dept_name,
	um_user_mail,
	um_user_phone,
	um_user_xerox,
	um_user_pass,
	um_flg_use,
	um_corp_cd,
	um_corp_name,
	um_refer_order_table,
	um_publisher_flg,
	um_publisher_id,
	um_pw_expiring_date,
	um_pw_change_flg,
	um_option_flg_return
)
VALUES
` + strValue
	//
	err = this.executeTrans(queryStr, args)
	//
	return randomPass, exterror.WrapExtError(err)
}

// Update
func (this *UserMasterAPI_Model) UpdateUser(parameter []byte) error {
	//
	newUser := ModelItems.UserItem{}
	err := json.Unmarshal(parameter, &newUser)
	if err != nil {
		return exterror.WrapExtError(err)
	}
	//
	newUser.UpdateDate = Common.CurrentDatetime()
	// check newUser.OptionFlgReturn
	args := []interface{}{
		newUser.UpdateDate,
		newUser.UserName,
		newUser.FlgAuth,
		newUser.ShopChainCd,
		newUser.FranchiseCd,
		newUser.FranchiseGroupCd,
		newUser.ServerName,
		newUser.ShopCd,
		newUser.ShopName,
		newUser.FlgMenuGroup,
		newUser.DeptCd,
		newUser.DeptName,
		newUser.UserMail,
		newUser.UserPhone,
		newUser.UserXerox,
		newUser.FlgUse,
		newUser.CorpCd,
		newUser.CorpName,
		newUser.ReferOrderTable,
		newUser.PublisherFlg,
		newUser.PublisherId,
	}
	setValueOptionFlgReturn := ""
	if newUser.OptionFlgReturn != "" {
		args = append(args, newUser.OptionFlgReturn)
		setValueOptionFlgReturn = `
			,um_option_flg_return  = ? 
		`
	}
	args = append(args, newUser.UserID)
	//
	queryStr := `
UPDATE :um_tb_name
SET
	um_update_date        = ?,
	um_user_name          = ?,
	um_flg_auth           = ?,
	um_shop_chain_cd      = ?,
	um_franchise_cd       = ?,
	um_franchise_group_cd = ?,
	um_server_name        = ?,
	um_shop_cd            = ?,
	um_shop_name          = ?,
	um_flg_menu_group     = ?,
	um_dept_cd            = ?,
	um_dept_name          = ?,
	um_user_mail          = ?,
	um_user_phone         = ?,
	um_user_xerox         = ?,
	um_flg_use            = ?,
	um_corp_cd            = ?,
	um_corp_name          = ?,
	um_refer_order_table  = ?,
	um_publisher_flg      = ?,
	um_publisher_id       = ?
	` + setValueOptionFlgReturn + `
WHERE TRUE
	AND um_user_ID = ?
`
	//
	err = this.executeTrans(queryStr, args)
	//
	return exterror.WrapExtError(err)
}

// Delete
func (this *UserMasterAPI_Model) DeleteUser(user_id string) error {
	//
	args := []interface{}{
		Common.CurrentDatetime(),
		user_id,
	}
	//
	queryStr := `
UPDATE :um_tb_name
SET
	um_update_date = ?,
	um_flg_use     = '1'
WHERE TRUE
	AND um_user_ID = ?
`
	//
	err := this.executeTrans(queryStr, args)
	//
	return exterror.WrapExtError(err)
}

// Reset password
func (this *UserMasterAPI_Model) ResetPassword(user_id string) (string, error) {
	//
	randomPass := Common.RandNumber(_PASSWORD_LENGTH_DEFAULT)
	encryptPass := Common.GeneratePass(user_id, randomPass)
	expiringDate := Common.DatetimePastWithFormat(1, Common.DATE_FORMAT_ZERO_TIME)
	//
	args := []interface{}{
		Common.CurrentDatetime(),
		encryptPass,
		expiringDate,
		user_id,
	}
	//
	queryStr := `
UPDATE :um_tb_name
SET
	um_update_date   = ?,
	um_user_pass     = ?,
	um_flg_use       = '0',
	um_pw_change_flg = '0',
	um_pw_expiring_date = ?
WHERE TRUE
	AND um_user_ID = ?
`
	//
	err := this.executeTrans(queryStr, args)
	//
	return randomPass, exterror.WrapExtError(err)
}
