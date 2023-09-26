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

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *categoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) Create(req *models.CreateCategory) (string, error) {
	var (
		id = uuid.NewString()
	)

	query := `
				INSERT INTO "category"(
					"id",
					"name",
					"created_at")
				VALUES ($1, $2, NOW())`

	if req.ParentId != "" {
		query = `
				INSERT INTO "category"(
					"id",
					"name",
					"parent_id",
					"created_at")
				VALUES ($1, $2, $3, NOW())`
	}

	if req.ParentId != "" {
		_, err := r.db.Exec(context.Background(), query,
			id,
			req.Name,
			req.ParentId,
		)

		if err != nil {
			return "", err
		}
	} else {
		_, err := r.db.Exec(context.Background(), query,
			id,
			req.Name,
		)

		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (r *categoryRepo) GetByID(req *models.CategoryPrimaryKey) (*models.Category, error) {

	var (
		id        sql.NullString
		name      sql.NullString
		parent_id sql.NullString
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query := `
		SELECT
			"id", 
			"name",
			"parent_id",
			"created_at",
			"updated_at" 
		FROM "category"
		WHERE id = $1
	`

	err := r.db.QueryRow(context.Background(), query, req.Id).Scan(
		&id,
		&name,
		&parent_id,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Category{
		Id:        id.String,
		Name:      name.String,
		ParentId:  parent_id.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}, nil
}

func (r *categoryRepo) GetList(req *models.CategoryGetListRequest) (*models.CategoryGetListResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.CategoryGetListResponse{}

	resp.Categories = make([]*models.Category, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id", 
				"name",
				"parent_id",
				"created_at",
				"updated_at" 
			FROM "category"
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
			id        sql.NullString
			name      sql.NullString
			parent_id sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&parent_id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Categories = append(resp.Categories, &models.Category{
			Id:        id.String,
			Name:      name.String,
			ParentId:  parent_id.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}
	return resp, nil
}
func (r *categoryRepo) Update(req *models.UpdateCategory) (string, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"category"
		SET
			"name" = :name,
			"updated_at" = NOW()
	`

	params = map[string]interface{}{
		"id":   req.Id,
		"name": req.Name,
	}

	if req.ParentId != "" {
		query += `,
			"parent_id" = :parent_id
		`
		params["parent_id"] = req.ParentId
	}

	query += `
		WHERE id = :id
	`

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("category with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (r *categoryRepo) Delete(req *models.CategoryPrimaryKey) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx, "DELETE FROM category WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category with ID %s not found", req.Id)

	}

	return nil
}
