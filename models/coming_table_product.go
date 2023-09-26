package models

type ComingTableProductPrimaryKey struct {
	Id string `json:"id"`
}

type ComingTableProductBarcode struct {
	Barcode       string `json:"barcode"`
	ComingTableId string `json:"coming_table_id"`
}

type CreateComingTableProductCount struct {
	Count int `json:"count"`
}

type CreateComingTableProduct struct {
	CategoryId     string  `json:"category_id"`
	ProductName    string  `json:"name"`
	ProductPrice   float64 `json:"price"`
	ProductBarcode string  `json:"barcode"`
	Count          int     `json:"count"`
	TotalPrice     float64 `json:"total_price"`
	ComingTableId  string  `json:"coming_table_id"`
}

type ComingTableProduct struct {
	Id             string  `json:"id"`
	CategoryId     string  `json:"category_id"`
	ProductName    string  `json:"name"`
	ProductPrice   float64 `json:"price"`
	ProductBarcode string  `json:"barcode"`
	Count          int     `json:"count"`
	TotalPrice     float64 `json:"total_price"`
	ComingTableId  string  `json:"coming_table_id"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type UpdateComingTableProduct struct {
	Id             string  `json:"id"`
	CategoryId     string  `json:"category_id"`
	ProductName    string  `json:"name"`
	ProductPrice   float64 `json:"price"`
	ProductBarcode string  `json:"barcode"`
	Count          int     `json:"count"`
	TotalPrice     float64 `json:"total_price"`
	ComingTableId  string  `json:"coming_table_id"`
}

type ComingTableProductGetListRequest struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	CategoryId     string `json:"category_id"`
	ProductBarcode string `json:"barcode"`
}

type ComingTableProductGetListResponse struct {
	Count               int                   `json:"count"`
	ComingTableProducts []*ComingTableProduct `json:"products"`
}
