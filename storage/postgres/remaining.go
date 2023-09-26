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

type remainingRepo struct {
	db *pgxpool.Pool
}

func NewRemainingRepo(db *pgxpool.Pool) *remainingRepo {
	return &remainingRepo{
		db: db,
	}
}

func (r *remainingRepo) Create(req *models.CreateRemaining) (string, error) {
	var (
		id    = uuid.NewString()
		query string
	)

	query = `
		INSERT INTO "remaining"(
			"id", 
			"branch_id",
			"category_id",
			"name",
			"price",
			"barcode",
			"count",
			"total_price",
			"created_at" )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.BranchId,
		req.CategoryId,
		req.Name,
		req.Price,
		req.Barcode,
		req.Count,
		req.TotalPrice,
	)

	if err != nil {
		return "", err
	}

	return id, nil

}

func (r *remainingRepo) GetByID(req *models.RemainingPrimaryKey) (*models.Remaining, error) {
	var updatedAt sql.NullString
	var createdAt sql.NullString

	query := `
		SELECT
			"id", 
			"branch_id",
			"category_id",
			"name",
			"price",
			"barcode",
			"count",
			"total_price",
			"created_at",
			"updated_at" 
		FROM "remaining"
		WHERE id = $1
	`
	remaining := models.Remaining{}
	err := r.db.QueryRow(context.Background(), query, req.Id).Scan(
		&remaining.Id,
		&remaining.BranchId,
		&remaining.CategoryId,
		&remaining.Name,
		&remaining.Price,
		&remaining.Barcode,
		&remaining.Count,
		&remaining.TotalPrice,
		&createdAt,
		&updatedAt,
	)
	remaining.CreatedAt = createdAt.String
	remaining.UpdatedAt = updatedAt.String
	if err != nil {
		return nil, err
	}

	return &remaining, nil
}

func (r *remainingRepo) GetList(req *models.RemainingGetListRequest) (*models.RemainingGetListResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.RemainingGetListResponse{}
	resp.Remainings = make([]*models.Remaining, 0)

	filter := " WHERE true "
	query := `
		SELECT
			COUNT(*) OVER(),
			"id", 
			"branch_id",
			"category_id",
			"name",
			"price",
			"barcode",
			"count",
			"total_price",
			"created_at",
			"updated_at" 
		FROM "remaining"
		`
	if req.BranchId != "" {
		filter += ` AND ("branch_id" = :branch_id) `
		params["branch_id"] = req.BranchId
	}

	if req.Barcode != "" {
		filter += ` AND ("barcode" = :barcode) `
		params["barcode"] = req.Barcode
	}

	if req.CategoryId != "" {
		filter += ` AND ("category_id" = :category_id) `
		params["category_id"] = req.CategoryId
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
		var updatedAt sql.NullString
		var createdAt sql.NullString

		var remaining models.Remaining
		err := rows.Scan(
			&resp.Count,
			&remaining.Id,
			&remaining.BranchId,
			&remaining.CategoryId,
			&remaining.Name,
			&remaining.Price,
			&remaining.Barcode,
			&remaining.Count,
			&remaining.TotalPrice,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		remaining.CreatedAt = createdAt.String
		remaining.UpdatedAt = updatedAt.String

		resp.Remainings = append(resp.Remainings, &remaining)
	}
	return resp, nil
}

func (r *remainingRepo) Update(req *models.UpdateRemaining) (string, error) {

	query := `
		UPDATE
			"remaining"
		SET
			"branch_id" = $1,
			"category_id" = $2,
			"name" = $3,
			"price" = $4,
			"barcode" =$5,
			"count" = $6,
			"total_price" =$7,
			"updated_at" = NOW()
		WHERE id = $8
	`

	result, err := r.db.Exec(context.Background(), query,
		req.BranchId,
		req.CategoryId,
		req.Name,
		req.Price,
		req.Barcode,
		req.Count,
		req.TotalPrice,
		req.Id,
	)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("remaining with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (r *remainingRepo) Delete(req *models.RemainingPrimaryKey) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx, "DELETE FROM remaining WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("remaining with ID %s not found", req.Id)

	}

	return nil
}

// check raming by branch id and barcode
func (r *remainingRepo) CheckRemaing(req *models.CheckingRemaining) (string, error) {
	var in sql.NullString
	var params map[string]interface{}

	query := `
		SELECT
			"id"
		FROM "remaining"
		WHERE "branch_id" = :branch_id AND "barcode" = :barcode
	`

	params = map[string]interface{}{
		"branch_id": req.BranchId,
		"barcode":   req.Barcode,
	}
	queryN, args := helper.ReplaceQueryParams(query, params)

	err := r.db.QueryRow(context.Background(), queryN, args...).Scan(
		&in,
	)
	if err != nil {
		return in.String, err
	}

	return in.String, nil
}

func (r *remainingRepo) UpdateExists(req *models.UpdateRemaining) (string, error) {

	query := `
		UPDATE
			"remaining"
		SET
			"branch_id" = $1,
			"category_id" = $2,
			"name" = $3,
			"price" = $4,
			"barcode" =$5,
			"count" = "count" + $6,
			"total_price" = "total_price" + $7,
			"updated_at" = NOW()
		WHERE id = $8
	`

	result, err := r.db.Exec(context.Background(), query,
		req.BranchId,
		req.CategoryId,
		req.Name,
		req.Price,
		req.Barcode,
		req.Count,
		req.TotalPrice,
		req.Id,
	)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("remaining with ID %s not found", req.Id)
	}

	return req.Id, nil
}
