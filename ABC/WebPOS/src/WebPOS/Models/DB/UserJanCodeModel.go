package Models

import (
	"WebPOS/Common"
	"WebPOS/Models/ModelItems"
	"database/sql"
	"github.com/goframework/gf/db"
	"github.com/goframework/gf/exterror"
)

type UserJanCodeModel struct {
	DB *sql.DB
	trans *sql.Tx
}


func NewUserJanCodeModel(DB *sql.DB) UserJanCodeModel {
	model := UserJanCodeModel{}
	model.DB = DB

	return model
}

func (this *UserJanCodeModel) GetListUserJanByUser(userId string) (*[]ModelItems.UserJanCodeItem, error) {

	sqlString := `
SELECT
	ujc_create_datetime,
	ujc_update_datetime,
	ujc_user_id,
	ujc_jan_code,
	ujc_product_name,
	ujc_maker_name,
	ujc_author_name,
	ujc_selling_date,
	ujc_list_price,
	ujc_user_jan_inf,
	ujc_priority_number
FROM
	user_jan
WHERE ujc_user_id = ?
ORDER BY ujc_priority_number DESC, ujc_update_datetime DESC
`
	listUserJan := []ModelItems.UserJanCodeItem{}
	rows, err := this.DB.Query(sqlString, userId)
	if err != nil {
		return &listUserJan, exterror.WrapExtError(err)
	}
	defer rows.Close()
	for rows.Next() {
		newItem := ModelItems.UserJanCodeItem{}
		err := db.SqlScanStruct(rows, &newItem)
		if err != nil {
			return &listUserJan, exterror.WrapExtError(err)
		}
		listUserJan = append(listUserJan, newItem)
	}
	return &listUserJan, nil
}

func (this *UserJanCodeModel) InsertUpdateUserJan(args []interface{}) error {
	query := `
INSERT INTO user_jan (
	ujc_create_datetime,
	ujc_update_datetime,
	ujc_user_id,
	ujc_jan_code,
	ujc_product_name,
	ujc_maker_name,
	ujc_author_name,
	ujc_selling_date,
	ujc_list_price,
	ujc_user_jan_inf,
	ujc_priority_number)
VALUES (now(), now(), ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	user_jan.ujc_update_datetime = now(),
	user_jan.ujc_priority_number = user_jan.ujc_priority_number + 1,
	user_jan.ujc_product_name = IF(VALUES(ujc_product_name) <> "", VALUES(ujc_product_name), user_jan.ujc_product_name),
	user_jan.ujc_maker_name = IF(VALUES(ujc_maker_name) <> "", VALUES(ujc_maker_name), user_jan.ujc_maker_name),
	user_jan.ujc_author_name = IF(VALUES(ujc_author_name) <> "", VALUES(ujc_author_name), user_jan.ujc_author_name),
	user_jan.ujc_selling_date = IF(VALUES(ujc_selling_date) <> "", VALUES(ujc_selling_date), user_jan.ujc_selling_date),
	user_jan.ujc_list_price = IF(VALUES(ujc_list_price) <> "", VALUES(ujc_list_price), user_jan.ujc_list_price)
`
	_, err := this.Exec(query, args)
	return exterror.WrapExtError(err)
}

func (this *UserJanCodeModel) InsertUpdateUserJanSearchGoods(args []interface{}) error {
	query := `
INSERT INTO user_jan (
	ujc_create_datetime,
	ujc_update_datetime,
	ujc_user_id,
	ujc_jan_code,
	ujc_product_name,
	ujc_maker_name,
	ujc_author_name,
	ujc_selling_date,
	ujc_list_price,
	ujc_user_jan_inf,
	ujc_priority_number)
VALUES (now(), now(), ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	user_jan.ujc_update_datetime = now(),
	user_jan.ujc_priority_number = user_jan.ujc_priority_number + 100,
	user_jan.ujc_product_name = VALUES(ujc_product_name),
	user_jan.ujc_maker_name = VALUES(ujc_maker_name),
	user_jan.ujc_author_name = VALUES(ujc_author_name),
	user_jan.ujc_selling_date = VALUES(ujc_selling_date),
	user_jan.ujc_list_price = VALUES(ujc_list_price)
`
	_, err := this.DB.Exec(query, args...)
	return exterror.WrapExtError(err)
}

func (this *UserJanCodeModel) UpdateInfoListUserJan(args []interface{}) error {
	sqlQuery := `
UPDATE user_jan
SET
	ujc_update_datetime = now(),
	ujc_priority_number = ?,
	ujc_user_jan_inf = ?
WHERE
	ujc_user_id = ?
	AND ujc_jan_code = ?
	`
	_, err := this.Exec(sqlQuery, args)

	return exterror.WrapExtError(err)
}

func (this *UserJanCodeModel) DeleteListUserJan(userId string, janList []string) error {
	args := []interface{}{}
	sqlQuery := `
DELETE
FROM user_jan
WHERE
	ujc_user_id = ?
	AND ujc_jan_code IN(` + Common.SQLPara(janList) + `)
	`
	args = append(args, userId)
	args = append(args, Common.ToInterfaceArray(janList)...)
	_, err := this.DB.Exec(sqlQuery, args...)

	return exterror.WrapExtError(err)
}

func (this *UserJanCodeModel) DeleteAutoUserJan(userId string) error {

	sqlQuery := `
DELETE
FROM user_jan
WHERE
	ujc_user_id = ?
	AND CAST(ujc_update_datetime AS DATE) <= ADDDATE(CAST(NOW() AS DATE), INTERVAL -1 MONTH)
	`
	_, err := this.DB.Exec(sqlQuery, userId)

	return exterror.WrapExtError(err)
}

func (this *UserJanCodeModel) BeginTrans() error {
	if this.DB != nil && this.trans == nil {
		var err error
		if this.trans , err = this.DB.Begin(); err != nil {
			return exterror.WrapExtError(err)
		}
	}
	return nil
}

func (this *UserJanCodeModel) FinishTrans(isErr *error) error {
	if this.trans != nil {
		if *isErr != nil {
			if err := this.trans.Rollback(); err != nil {
				return exterror.WrapExtError(err)
			}
		}else{
			if err := this.trans.Commit(); err != nil {
				return exterror.WrapExtError(err)
			}
		}
	}
	return nil
}

func (this *UserJanCodeModel) Exec(query string, argument []interface{}) (sql.Result, error) {

	if argument == nil || len(argument) == 0 {
		return this.trans.Exec(query)
	} else {
		return this.trans.Exec(query, argument...)
	}
}