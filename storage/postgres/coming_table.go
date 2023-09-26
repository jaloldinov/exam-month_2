package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"market/models"
	"market/pkg/helper"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type comingTableRepo struct {
	db *pgxpool.Pool
}

func NewComingTableRepo(db *pgxpool.Pool) *comingTableRepo {
	return &comingTableRepo{
		db: db,
	}
}

func (r *comingTableRepo) Create(req *models.CreateComingTable) (string, error) {
	var (
		id = uuid.NewString()
	)

	query := `
				INSERT INTO "coming_table"(
					"id",
					"coming_id",
					"branch_id",
					"date_time",
					"created_at")
				VALUES ($1, $2, $3, $4, NOW())`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.ComingId,
		req.BranchId,
		req.DateTime,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *comingTableRepo) GetByID(req *models.ComingTablePrimaryKey) (*models.ComingTable, error) {

	var (
		id         sql.NullString
		coming_id  sql.NullString
		branch_id  sql.NullString
		date_time  sql.NullTime
		status     sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	query := `
		SELECT
			"id", 
			"coming_id",
			"branch_id",
			"date_time",
			"status",
			"created_at",
			"updated_at" 
		FROM "coming_table"
		WHERE id = $1
	`

	err := r.db.QueryRow(context.Background(), query, req.Id).Scan(
		&id,
		&coming_id,
		&branch_id,
		&date_time,
		&status,
		&created_at,
		&updated_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.ComingTable{
		Id:        id.String,
		ComingId:  coming_id.String,
		BranchId:  branch_id.String,
		DateTime:  date_time.Time.Format(time.DateTime),
		Status:    status.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}, nil
}

func (r *comingTableRepo) GetList(req *models.ComingTableGetListRequest) (*models.ComingTableGetListResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.ComingTableGetListResponse{}

	resp.ComingTables = make([]*models.ComingTable, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"coming_id",
				"branch_id",
				"date_time",
				"status",
				"created_at",
				"updated_at" 
			FROM "coming_table"
		`
	if req.ComingId != "" {
		filter += ` AND ("coming_id" = :coming_id)`
		params["coming_id"] = req.ComingId
	}

	if req.BranchId != "" {
		filter += ` AND ("branch_id" = :branch_id)`
		params["branch_id"] = req.BranchId
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
			id         sql.NullString
			coming_id  sql.NullString
			branch_id  sql.NullString
			date_time  sql.NullTime
			status     sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&coming_id,
			&branch_id,
			&date_time,
			&status,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		resp.ComingTables = append(resp.ComingTables, &models.ComingTable{
			Id:        id.String,
			ComingId:  coming_id.String,
			BranchId:  branch_id.String,
			DateTime:  date_time.Time.Format(time.DateTime),
			Status:    status.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}
	return resp, nil
}

func (r *comingTableRepo) Update(req *models.UpdateComingTable) (string, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"coming_table"
		SET
				"coming_id" = :coming_id
				"branch_id" = :branch_id
				"date_time" = :date_time
				"updated_at" = NOW()
				WHERE id = :id
	`

	params = map[string]interface{}{
		"id":        req.Id,
		"coming_id": req.ComingId,
		"branch_id": req.BranchId,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("coming_table with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (r *comingTableRepo) Delete(req *models.ComingTablePrimaryKey) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx, "DELETE FROM coming_table WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("coming_table with ID %s not found", req.Id)

	}

	return nil
}

func (r *comingTableRepo) UpdateStatus(req *models.ComingTablePrimaryKey) (string, error) {

	query := `
		UPDATE
			"coming_table"
		SET
				"status" = $1,
				"updated_at" = NOW()
				WHERE id = $2
	`

	result, err := r.db.Exec(context.Background(), query, "finished", req.Id)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("coming_table with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (r *comingTableRepo) GetStatus(req *models.ComingTablePrimaryKey) (string, error) {
	var status sql.NullString
	var branch_id sql.NullString

	uuidValue, err := uuid.Parse(req.Id)
	if err != nil {
		return "", fmt.Errorf("invalid UUID format: %v", err)
	}

	query := `
		SELECT
			"status",
			"branch_id"
		FROM "coming_table"
		WHERE "id" = $1::uuid
	`

	err = r.db.QueryRow(context.Background(), query, uuidValue).Scan(
		&status,
		&branch_id,
	)
	if err != nil {
		return "", err
	}

	if status.String == "finished" {
		return "", fmt.Errorf("coming table already finished")
	}

	return branch_id.String, nil
}
