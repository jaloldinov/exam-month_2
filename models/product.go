package models

type ProductPrimaryKey struct {
	Id string `json:"id"`
}

type ProductBarcodeRequest struct {
	Barcode string `json:"barcode"`
}

type ProductBarcodeResponse struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryId string  `json:"category_id"`
}

type CreateProduct struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Barcode    string  `json:"barcode"`
	CategoryId string  `json:"category_id"`
}

type Product struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Barcode    string  `json:"barcode"`
	CategoryId string  `json:"category_id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type UpdateProduct struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Barcode    string  `json:"barcode"`
	CategoryId string  `json:"category_id"`
}

type ProductGetListRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
}

type ProductGetListResponse struct {
	Count    int        `json:"count"`
	Products []*Product `json:"products"`
}
