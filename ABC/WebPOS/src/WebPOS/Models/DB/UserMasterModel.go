package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
	"strings"
)

type UserMasterModel struct {
	DB *sql.DB
}

func (this *UserMasterModel) Login(user, pass string) (bool, error) {
	// Login check
	query := `
SELECT
	IF(COUNT(*), 'true', 'false')
FROM
	ao_user_master
WHERE
	um_user_ID = ?
	AND um_user_pass = ?
 	AND um_flg_use = '0'
	`

	args := []interface{}{}
	args = append(args, user)
	args = append(args, Common.GeneratePass(user, pass))

	var isAuthenticated bool
	err := this.DB.QueryRow(query, args...).Scan(&isAuthenticated)
	if err != nil {
		return false, exterror.WrapExtError(err)
	}

	// Update last login
	if isAuthenticated {
		query = `
			UPDATE ao_user_master
			   SET um_latest_login_time = ?, um_update_date = ?
			 WHERE um_user_ID = ?
		`
		this.DB.Exec(query, Common.CurrentDateTime(), Common.CurrentDate(), user)
	}
	return isAuthenticated, nil
}

func (this *UserMasterModel) GetUserInfoById(userId string) (*ModelItems.UserItem, error) {
	query := `
		SELECT
				um_shop_cd
				,IFNULL(um_shop_name,'')
				,um_server_name
				,um_shop_chain_cd
				,um_flg_auth
				,um_user_ID
				,um_user_name
				,IFNULL(um_franchise_cd,'')
				,um_franchise_group_cd
				,um_flg_menu_group
				,IFNULL(um_user_mail,'')
				,IFNULL(um_user_phone,'')
				,IFNULL(um_user_xerox,'')
				,IFNULL(um_user_pass,'')
				,IFNULL(um_dept_name,'')
				,IFNULL(um_corp_cd,'')
				,IFNULL(um_corp_name,'')
				,IFNULL(um_refer_order_table,'')
				,um_pw_expiring_date
				,um_option_flg_return
				,IFNULL(um_dept_cd,'')
				,IFNULL(TRIM(um_publisher_id), '')
		FROM ao_user_master
		WHERE um_user_ID = ?
	`
	user := ModelItems.UserItem{}
	err := this.DB.QueryRow(query, userId).Scan(
		&user.ShopCd,
		&user.ShopName,
		&user.ServerName,
		&user.ShopChainCd,
		&user.FlgAuth,
		&user.UserID,
		&user.UserName,
		&user.FranchiseCd,
		&user.FranchiseGroupCd,
		&user.FlgMenuGroup,
		&user.UserMail,
		&user.UserPhone,
		&user.UserXerox,
		&user.UserPass,
		&user.DeptName,
		&user.CorpCd,
		&user.CorpName,
		&user.ReferOrderTable,
		&user.PwExpiringDate,
		&user.OptionFlgReturn,
		&user.DeptCd,
		&user.PublisherId,
	)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	return &user, nil
}

type MenuGroup struct {
	Menu_ID   string
	Menu_Name string
	Menu_Url  string
}

//func (this *UserMasterModel) GetMenuPathByGroup(flag_group string) (map[string]bool, error) {
//	query := `
//SELECT
//	IF(mcm_menu_ID like '/%', mcm_menu_ID, CONCAT('/', mcm_menu_ID)) menu_path
//FROM
//	menu_control_master
//WHERE mcm_flg_menu_group = ?
//`
//	mapMenu := map[string]bool{}
//	rows, err := this.DB.Query(query, flag_group)
//	if err != nil {
//		return nil, exterror.WrapExtError(err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		menu := ""
//		err = rows.Scan(&menu)
//		mapMenu[menu] = true
//	}
//
//	return mapMenu, nil
//}

func (this *UserMasterModel) GetMenuPathByGroup() (map[string]bool, error) {
	query := `
SELECT
	IF(menu_id like '/%', menu_id, CONCAT('/', menu_id)) menu_path
FROM
	menu_master
`
	mapMenu := map[string]bool{}
	rows, err := this.DB.Query(query)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	for rows.Next() {
		menu := ""
		err = rows.Scan(&menu)
		mapMenu[menu] = true
	}

	return mapMenu, nil
}

type ListUser struct {
	Um_User_ID            string
	Um_User_Name          string
	Um_Shop_Chain_Cd      string
	Um_Shop_Cd            string
	Um_Shop_Name          string
	Um_Flg_Auth           string
	Um_Franchise_Cd       string
	Um_Franchise_Group_Cd string
	Um_Server_Name        string
	Um_Flg_Menu_Group     string
	Um_Latest_Login_Time  string
	Um_Corp_Cd            string
	Um_Corp_Name          string
	Um_Refer_Order_Table  string
	Um_Dept_Cd            string
	Um_Dept_Name          string
	Um_User_Mail          string
	Um_User_Phone         string
	Um_User_Xerox         string
	Um_User_Pass          string
	Um_Flg_Use            string
}

func (this *UserMasterModel) GetListUser(arrShop []string, user_name string, flg_auth string) ([]ListUser, error) {

	var args []interface{}
	user_name = "%" + user_name + "%"
	args = append(args, user_name)

	sWhereShop := ""
	if len(arrShop) > 0 {
		for _, s := range arrShop {
			args = append(args, s)
		}
		sWhereShop = " AND um_shop_cd IN (?" + strings.Repeat(",?", len(arrShop)-1) + ")"
	}

	sWhereAuth := ""
	if flg_auth == "0" {
		sWhereAuth = " AND um_flg_auth = '0' "
	} else if flg_auth == "1" {
		sWhereAuth = " AND um_flg_auth = '1' "
	} else {
		sWhereAuth = ""
	}

	query := `
		SELECT  IFNULL(um_user_id,'')
				,IFNULL(um_user_name,'')
				,IFNULL(um_shop_chain_cd,'')
				,IFNULL(um_shop_cd,'')
				,IFNULL(um_shop_name,'')
				,IFNULL(um_flg_auth,'')
				,IFNULL(um_franchise_cd,'')
				,IFNULL(um_franchise_group_cd,'')
				,IFNULL(um_server_name,'')
				,IFNULL(um_flg_menu_group,'')
				,IFNULL(um_latest_login_time,'')
				,IFNULL(um_corp_cd,'')
				,IFNULL(um_corp_name,'')
				,IFNULL(um_refer_order_table,'')
				,IFNULL(um_dept_cd,'')
				,IFNULL(um_dept_name,'')
				,IFNULL(um_user_mail,'')
				,IFNULL(um_user_phone,'')
				,IFNULL(um_user_xerox,'')
				,IFNULL(um_user_pass,'')
				,IFNULL(um_flg_use,'')
		 FROM ao_user_master
		 WHERE um_user_name LIKE ?
	` + sWhereShop + sWhereAuth
	var list_user []ListUser
	rows, err := this.DB.Query(query, args...)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		g := ListUser{}
		err = rows.Scan(
			&g.Um_User_ID,
			&g.Um_User_Name,
			&g.Um_Shop_Chain_Cd,
			&g.Um_Shop_Cd,
			&g.Um_Shop_Name,
			&g.Um_Flg_Auth,
			&g.Um_Franchise_Cd,
			&g.Um_Franchise_Group_Cd,
			&g.Um_Server_Name,
			&g.Um_Flg_Menu_Group,
			&g.Um_Latest_Login_Time,
			&g.Um_Corp_Cd,
			&g.Um_Corp_Name,
			&g.Um_Refer_Order_Table,
			&g.Um_Dept_Cd,
			&g.Um_Dept_Name,
			&g.Um_User_Mail,
			&g.Um_User_Phone,
			&g.Um_User_Xerox,
			&g.Um_User_Pass,
			&g.Um_Flg_Use)
		if g.Um_Shop_Cd == "" {
			g.Um_Shop_Cd = "-1"
		}
		if g.Um_Franchise_Cd == "" {
			g.Um_Franchise_Cd = "-1"
		}
		list_user = append(list_user, g)
	}

	return list_user, nil
}

func (this *UserMasterModel) UpdateInsertUser(list_user []ListUser) error {
	date_update := Common.CurrentDate()

	len_para := 0
	args := []interface{}{}
	for _, list := range list_user {
		pass, err := this.CheckPass(list.Um_User_ID, list.Um_User_Pass)
		if err != nil {
			return exterror.WrapExtError(err)
		}

		args = append(args, list.Um_Shop_Chain_Cd)
		args = append(args, list.Um_User_ID)
		args = append(args, list.Um_User_Name)
		args = append(args, list.Um_Shop_Cd)
		args = append(args, list.Um_Shop_Name)
		args = append(args, list.Um_Flg_Auth)
		args = append(args, list.Um_Franchise_Cd)
		args = append(args, list.Um_Franchise_Group_Cd)
		args = append(args, list.Um_Server_Name)
		args = append(args, list.Um_Flg_Menu_Group)
		args = append(args, list.Um_Corp_Cd)
		args = append(args, list.Um_Corp_Name)
		args = append(args, list.Um_Refer_Order_Table)
		args = append(args, list.Um_Dept_Cd)
		args = append(args, list.Um_Dept_Name)
		args = append(args, list.Um_User_Mail)
		args = append(args, list.Um_User_Phone)
		args = append(args, list.Um_User_Xerox)
		args = append(args, pass)
		args = append(args, list.Um_Flg_Use)
		args = append(args, date_update)
		args = append(args, date_update)
		len_para = len_para + 1
	}
	strValue := strings.Repeat("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),", len_para)
	strValue = strings.TrimRight(strValue, ",")

	var sqlString string
	sqlString = `	INSERT INTO ao_user_master (
								 um_shop_chain_cd
								,um_user_ID
								,um_user_name
								,um_shop_cd
								,um_shop_name
								,um_flg_auth
								,um_franchise_cd
								,um_franchise_group_cd
								,um_server_name
								,um_flg_menu_group
								,um_corp_cd
								,um_corp_name
								,um_refer_order_table
								,um_dept_cd
								,um_dept_name
								,um_user_mail
								,um_user_phone
								,um_user_xerox
								,um_user_pass
								,um_flg_use
								,um_create_date
								,um_update_date)
					VALUES
					` + strValue + `
					ON DUPLICATE KEY UPDATE  um_shop_chain_cd			= VALUES(um_shop_chain_cd)
											,um_user_name	 			= VALUES(um_user_name)
											,um_shop_cd	 				= VALUES(um_shop_cd)
											,um_shop_name	 			= VALUES(um_shop_name)
											,um_flg_auth	 			= VALUES(um_flg_auth)
											,um_franchise_cd	 		= VALUES(um_franchise_cd)
											,um_franchise_group_cd	 	= VALUES(um_franchise_group_cd)
											,um_server_name	 			= VALUES(um_server_name)
											,um_flg_menu_group	 		= VALUES(um_flg_menu_group)
											,um_corp_cd	 				= VALUES(um_corp_cd)
											,um_corp_name	 			= VALUES(um_corp_name)
											,um_refer_order_table	 	= VALUES(um_refer_order_table)
											,um_dept_cd	 				= VALUES(um_dept_cd)
											,um_dept_name	 			= VALUES(um_dept_name)
											,um_user_mail	 			= VALUES(um_user_mail)
											,um_user_phone	 			= VALUES(um_user_phone)
											,um_user_xerox	 			= VALUES(um_user_xerox)
											,um_user_pass	 			= VALUES(um_user_pass)
											,um_flg_use	 				= VALUES(um_flg_use)
											,um_update_date				= VALUES(um_update_date)
	`
	_, err := this.DB.Exec(sqlString, args...)

	return exterror.WrapExtError(err)
}

func (this *UserMasterModel) CheckPass(user, pass string) (string, error) {
	// Login check
	query := `
		SELECT IF(COUNT(*), 'true', 'false')
		  FROM ao_user_master
		 WHERE um_user_ID = ?
		   AND um_user_pass = ?
	`
	args := []interface{}{}
	args = append(args, user)
	args = append(args, pass)

	var isAuthenticated bool
	err := this.DB.QueryRow(query, args...).Scan(&isAuthenticated)
	if err != nil {
		return pass, exterror.WrapExtError(err)
	}
	if !isAuthenticated {
		pass = Common.GeneratePass(user, pass)
	}
	return pass, nil
}

func (this *UserMasterModel) CheckUser(user string) (bool, error) {
	// Login check
	query := `
		SELECT IF(COUNT(*), 'true', 'false')
		  FROM ao_user_master
		 WHERE um_user_ID = ?
	`
	args := []interface{}{}
	args = append(args, user)

	var isAuthenticated bool
	err := this.DB.QueryRow(query, args...).Scan(&isAuthenticated)
	if err != nil {
		return false, exterror.WrapExtError(err)
	}
	return isAuthenticated, nil
}

// ユーザメンテナンス
//----------------------------------    パスワード変更     ----------------------------------//
func (this *UserMasterModel) ChangePass(user, pass string, pwExpiringDate int) error {
	query := `
		UPDATE ao_user_master
		SET
			um_user_pass = ?,
			um_update_date = ?,
			um_pw_expiring_date = ?
		WHERE um_user_ID = ?
	`
	newExpiringDate := Common.DatetimePastWithFormat(pwExpiringDate*(-1), Common.DATE_FORMAT_ZERO_TIME)
	args := []interface{}{
		Common.GeneratePass(user, pass),
		Common.CurrentDateTime(),
		newExpiringDate,
		user,
	}
	_, err := this.DB.Exec(query, args...)
	return exterror.WrapExtError(err)
}

func (this *UserMasterModel) CsvMultipleRegister(args []interface{}) (error, []string) {
	date := Common.CurrentDate()
	strTemplate := " ( '" + date + "', '" + date + "' , ?, ?, '', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ,"
	str := strings.Repeat(strTemplate, len(args)/20)
	str = strings.TrimRight(str, ",")
	if str == "" {
		return nil, []string{}
	}
	query := `
		INSERT IGNORE INTO ao_user_master (um_create_date,
										um_update_date,
										um_user_ID,
										um_user_name,
										um_shop_chain_cd,
										um_shop_cd,
										um_shop_name,
										um_flg_auth,
										um_server_name,
										um_franchise_cd,
										um_franchise_group_cd,
										um_flg_menu_group,
										um_latest_login_time,
										um_dept_cd,
										um_dept_name,
										um_user_mail,
										um_user_phone,
										um_user_xerox,
										um_user_pass,
										um_flg_use,
										um_corp_cd,
										um_corp_name,
										um_refer_order_table)
		VALUES ` + str + `;
	`
	_, err := this.DB.Exec(query, args...)
	if err != nil {
		return err, []string{}
	}
	return this.getWarnings()
}

func (this *UserMasterModel) getWarnings() (error, []string) {
	rs := []string{}
	query := `SHOW WARNINGS;`
	rows, err := this.DB.Query(query)

	defer rows.Close()
	for rows.Next() {
		stt := ""
		code := ""
		msg := ""
		err = rows.Scan(
			&stt, &code, &msg,
		)
		if code == KEY_DUPLICATE {
			rs = append(rs, msg)
		}
	}
	return err, rs
}

//▼▼▼▼▼▼▼▼▼
func (this *UserMasterModel) CheckUserLogin(user string) (bool, error) {
	query := `
		SELECT IF(COUNT(*), 'true', 'false')
		  FROM ao_user_master
		 WHERE um_user_ID = ?
		 AND um_flg_use = '0'
	`
	args := []interface{}{
		user,
	}
	var isAuthenticated bool
	err := this.DB.QueryRow(query, args...).Scan(&isAuthenticated)
	if err != nil {
		return false, exterror.WrapExtError(err)
	}
	return isAuthenticated, nil
}
func (this *UserMasterModel) CheckUserPassLogin(user string, pass string) (bool, error) {
	query := `
		SELECT IF(COUNT(*), 'true', 'false')
		  FROM ao_user_master
		 WHERE um_user_ID = ?
		   AND um_user_pass = ?
		   AND um_flg_use = '0'
	`
	args := []interface{}{
		user,
		Common.GeneratePass(user, pass),
	}

	var isAuthenticated bool
	err := this.DB.QueryRow(query, args...).Scan(&isAuthenticated)
	if err != nil {
		return false, exterror.WrapExtError(err)
	}
	return isAuthenticated, nil
}

//▲▲▲▲▲▲▲▲▲
type MenuGroupUser struct {
	Menu_Group string
}

func (this *UserMasterModel) ListMenuGroup(chain_cd, shop_cd string) []MenuGroupUser {
	var listMenuGroup []MenuGroupUser
	//var sqlString string = ""
	//
	//sqlString = `
	//			SELECT mcm_flg_menu_group FROM menu_control_master GROUP BY mcm_flg_menu_group
	//	`
	//rows, err := this.DB.Query(sqlString)
	//
	//if err != nil {
	//	return nil
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	g := MenuGroupUser{}
	//	err = rows.Scan(&g.Menu_Group)
	//	listMenuGroup = append(listMenuGroup, g)
	//}
	return listMenuGroup
}

func (this UserMasterModel) GetUserRole(userId string) (*ModelItems.UserRoleItem, error) {
	sql := `
SELECT
	um_user_ID,
	um_shop_cd,
	um_server_name,
	um_flg_auth,
	um_shop_chain_cd,
	um_franchise_cd,
	um_franchise_group_cd
FROM ao_user_master
WHERE um_user_ID = ?
	`
	rows, err := this.DB.Query(sql, userId)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}
	defer rows.Close()

	if rows.Next() {
		item := ModelItems.UserRoleItem{}
		err = db.SqlScanStruct(rows, &item)
		if err != nil {
			return nil, exterror.WrapExtError(err)
		}

		item.ListChainCD = strings.Split(item.ChainCD, ",")
		return &item, nil
	}

	return nil, nil
}

// Update 「パスワード間違い回数」 when login
func (this *UserMasterModel) UpdateChangePWFlag(userID string, isAutoIncrement bool) error {

	sql := ""
	if isAutoIncrement {
		sql = `
UPDATE ao_user_master
SET
	um_update_date = ?,
	um_pw_change_flg = um_pw_change_flg + 1
WHERE
	um_user_ID = ?
	`
	} else {
		sql = `
UPDATE ao_user_master
SET
	um_update_date = ?,
	um_pw_change_flg = 0
WHERE
	um_user_ID = ?
	`
	}
	_, err := this.DB.Exec(sql, Common.CurrentDate(), userID)

	return exterror.WrapExtError(err)
}

// Get 「パスワード有効期限」 by user ID
func (this *UserMasterModel) GetChangePWTime(user string) (int, error) {
	// Login check
	query := `
SELECT
	um_pw_change_flg
FROM
	ao_user_master
WHERE
	um_user_ID = ?
 	AND um_flg_use = '0'
	`

	var PwChangeFlg int
	err := this.DB.QueryRow(query, user).Scan(&PwChangeFlg)
	if err != nil {
		return -1, exterror.WrapExtError(err)
	}
	return PwChangeFlg, nil
}