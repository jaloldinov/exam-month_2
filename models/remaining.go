package models

type RemainingPrimaryKey struct {
	Id string `json:"id"`
}
type CheckingRemaining struct {
	BranchId string `json:"branch_id"`
	Barcode  string `json:"barcode"`
}
type CreateRemaining struct {
	BranchId   string  `json:"branch_id"`
	CategoryId string  `json:"category_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Barcode    string  `json:"barcode"`
	Count      int     `json:"count"`
	TotalPrice float64 `json:"total_price"`
}

type Remaining struct {
	Id         string  `json:"id"`
	BranchId   string  `json:"branch_id"`
	CategoryId string  `json:"category_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Barcode    string  `json:"barcode"`
	Count      int     `json:"count"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type UpdateRemaining struct {
	Id         string  `json:"id"`
	BranchId   string  `json:"branch_id"`
	CategoryId string  `json:"category_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Barcode    string  `json:"barcode"`
	Count      int     `json:"count"`
	TotalPrice float64 `json:"total_price"`
}

type UpdateRemainingSoft struct {
	BranchId   string  `json:"branch_id"`
	CategoryId string  `json:"category_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Barcode    string  `json:"barcode"`
	Count      int     `json:"count"`
}

type RemainingGetListRequest struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	BranchId   string `json:"branch_id"`
	Barcode    string `json:"barcode"`
	CategoryId string `json:"category_id"`
}

type RemainingGetListResponse struct {
	Count      int          `json:"count"`
	Remainings []*Remaining `json:"remainings"`
}
