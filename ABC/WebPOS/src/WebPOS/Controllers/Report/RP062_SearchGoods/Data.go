package RP062_SearchGoods

type Response struct {
	Status         string    `json:"status"`
	Message        string    `json:"message"`
	ProductList    []Product `json:"product_list"`
	TotalPageCount string    `json:"total_page_count"`
	TotalItemCount string    `json:"total_item_count"`
	DisplayPage    string    `json:"display_page"`
	DisplayCount   string    `json:"display_count"`
}
type Product struct {
	ProductName string `json:"product_name"`
	ISBN        string `json:"isbn"`
	Author      Info   `json:"author"`
	Publisher   Info   `json:"publisher"`
	Release     Info   `json:"release"`
	Price       Info   `json:"price"`
}
type Info struct {
	AuthorName    string `json:"author_name"`
	PublisherName string `json:"publisher_name"`
	ReleaseDate   string `json:"release_date"`
	UsualPrice    string `json:"usual_price"`
}
