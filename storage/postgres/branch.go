package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"market/models"
	"market/pkg/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *branchRepo {
	return &branchRepo{
		db: db,
	}
}

func (r *branchRepo) Create(req *models.CreateBranch) (string, error) {

	var (
		id    = uuid.NewString()
		query string
	)

	query = `
		INSERT INTO "branch"(
			"id", 
			"name",
			"address",
			"phone_number",
			"created_at" )
		VALUES ($1, $2, $3, $4, NOW())`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.Name,
		req.Address,
		req.PhoneNumber,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *branchRepo) GetByID(req *models.BranchPrimaryKey) (*models.Branch, error) {

	var (
		id          sql.NullString
		name        sql.NullString
		address     sql.NullString
		phoneNumber sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	query := `
		SELECT
			"id", 
			"name",
			"address",
			"phone_number",
			"created_at",
			"updated_at" 
		FROM "branch"
		WHERE id = $1
	`

	err := r.db.QueryRow(context.Background(), query, req.Id).Scan(
		&id,
		&name,
		&address,
		&phoneNumber,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Branch{
		Id:          id.String,
		Name:        name.String,
		Address:     address.String,
		PhoneNumber: phoneNumber.String,
		CreatedAt:   createdAt.String,
		UpdatedAt:   updatedAt.String,
	}, nil
}

func (r *branchRepo) GetList(req *models.BranchGetListRequest) (*models.BranchGetListResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.BranchGetListResponse{}

	resp.Branches = make([]*models.Branch, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"name",
				"address",
				"phone_number",
				"created_at",
				"updated_at" 
			FROM "branch"
		`
	if req.Search != "" {
		filter += ` AND "name" ILIKE '%' || :search || '%' `
		params["search"] = req.Search
	}

	offset := (req.Page - 1) * req.Limit
	params["limit"] = req.Limit
	params["offset"] = offset

	query = query + filter + " ORDER BY created_at DESC OFFSET :offset LIMIT :limit "
	rquery, pArr := helper.ReplaceQueryParams(query, params)

	rows, err := r.db.Query(context.Background(), rquery, pArr...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          sql.NullString
			name        sql.NullString
			address     sql.NullString
			phoneNumber sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&address,
			&phoneNumber,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Branches = append(resp.Branches, &models.Branch{
			Id:          id.String,
			Name:        name.String,
			Address:     address.String,
			PhoneNumber: phoneNumber.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}
	return resp, nil

}

func (r *branchRepo) Update(req *models.UpdateBranch) (string, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"branch"
		SET
			"name" = :name,
			"address" = :address,
			"phone_number" = :phone_number,
			"updated_at" = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"name":         req.Name,
		"address":      req.Address,
		"phone_number": req.PhoneNumber,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("branch with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (r *branchRepo) Delete(req *models.BranchPrimaryKey) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx, "DELETE FROM branch WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("branch with ID %s not found", req.Id)

	}

	return nil
}
