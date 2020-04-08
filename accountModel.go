package database

impport (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "123"
	dbname   = "mydb"
)

type accountModel struct {
	DB *sql.DB
}
func connectDB() {

	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Printf("Fail to openDB: %v \n", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Fail to conenct: %v \n", err)
		return
	}
	fmt.Println("Ping OK")

}

func queryData() {

	_sql := "SELECT user_name, full_name FROM public.account LIMIT 1;"

	row, err := db.Query(_sql)
	if err != nil {
		fmt.Printf("Fail to query: %v \n", err)
		return
	}

	var col1 string
	var col2 string
	for row.Next() {
		row.Scan(&col1, &col2)
		fmt.Printf("value Col1: %v \n", col1)
		fmt.Printf("value Col2: %v \n", col2)
	}

	fmt.Println("End !!!")

}
