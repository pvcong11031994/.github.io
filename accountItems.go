package database

type accountItems struct {
	userName string `sql:"user_name"`
	fullName string `sql:"full_name"`
}
