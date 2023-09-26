package models

type CategoryPrimaryKey struct {
	Id string `json:"id"`
}

type CreateCategory struct {
	Name     string `json:"name"`
	ParentId string `json:"parent_id"`
}

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ParentId  string `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateCategory struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ParentId string `json:"parent_id"`
}

type CategoryGetListRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type CategoryGetListResponse struct {
	Count      int         `json:"count"`
	Categories []*Category `json:"categories"`
}
