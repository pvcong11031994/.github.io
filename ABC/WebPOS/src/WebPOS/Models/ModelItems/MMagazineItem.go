package ModelItems

type MMagazineItem struct {
	MagazineCode string `sql:"magazine_code"`
	MakerCode    string `sql:"maker_code"`
	MagazineName string `sql:"magazine_name"`
}
