package models

type BranchPrimaryKey struct {
	Id string `json:"id"`
}

type CreateBranch struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type Branch struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateBranch struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type BranchGetListRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type BranchGetListResponse struct {
	Count    int       `json:"count"`
	Branches []*Branch `json:"branches"`
}
