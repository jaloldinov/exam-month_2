package models

type ComingTablePrimaryKey struct {
	Id string `json:"id"`
}

type CreateComingTable struct {
	ComingId string `json:"coming_id"`
	BranchId string `json:"branch_id"`
	DateTime string `json:"date_time"`
}

type ComingTable struct {
	Id        string `json:"id"`
	ComingId  string `json:"coming_id"`
	BranchId  string `json:"branch_id"`
	DateTime  string `json:"date_time"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ComingIdResponse struct {
	Status string `json:"status"`
}

type UpdateComingTable struct {
	Id       string `json:"id"`
	ComingId string `json:"coming_id"`
	BranchId string `json:"branch_id"`
	DateTime string `json:"date_time"`
}

type ComingTableGetListRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	BranchId string `json:"branch_id"`
	ComingId string `json:"coming_id"`
}

type ComingTableGetListResponse struct {
	Count        int            `json:"count"`
	ComingTables []*ComingTable `json:"products"`
}
