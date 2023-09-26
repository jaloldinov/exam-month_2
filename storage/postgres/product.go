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

type productRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(req *models.CreateProduct) (string, error) {
	var (
		id = uuid.NewString()
	)

	query := `
				INSERT INTO "product"(
					"id",
					"name",
					"price",
					"barcode",
					"category_id",
					"created_at")
				VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.Name,
		req.Price,
		req.Barcode,
		req.CategoryId,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *productRepo) GetByID(req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		id          sql.NullString
		name        sql.NullString
		price       sql.NullFloat64
		barcode     sql.NullString
		category_id sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	query := `
		SELECT
			"id",
			"name",
			"price",		
			"barcode",
			"category_id",
			"created_at",
			"updated_at" 
		FROM "product"
		WHERE id = $1
	`

	err := r.db.QueryRow(context.Background(), query, req.Id).Scan(
		&id,
		&name,
		&price,
		&barcode,
		&category_id,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Product{
		Id:         id.String,
		Name:       name.String,
		Price:      price.Float64,
		Barcode:    barcode.String,
		CategoryId: category_id.String,
		CreatedAt:  createdAt.String,
		UpdatedAt:  updatedAt.String,
	}, nil
}

func (r *productRepo) GetList(req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.ProductGetListResponse{}

	resp.Products = make([]*models.Product, 0)

	filter := " WHERE true "
	query := `
		SELECT
			COUNT(*) OVER(),
			"id",
			"name",
			"price",		
			"barcode",
			"category_id",
			"created_at",
			"updated_at" 
		FROM "product"
	`
	if req.Name != "" {
		filter += ` AND "name" ILIKE '%' || :name || '%' `
		params["name"] = req.Name
	}

	if req.Barcode != "" {
		filter += ` AND ("barcode" = :barcode) `
		params["barcode"] = req.Barcode
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
			price       sql.NullFloat64
			barcode     sql.NullString
			category_id sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&price,
			&barcode,
			&category_id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		resp.Products = append(resp.Products, &models.Product{
			Id:         id.String,
			Name:       name.String,
			Price:      price.Float64,
			Barcode:    barcode.String,
			CategoryId: category_id.String,
			CreatedAt:  createdAt.String,
			UpdatedAt:  updatedAt.String,
		})
	}
	return resp, nil
}

func (r *productRepo) Update(req *models.UpdateProduct) (string, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"product"
		SET
			"name" = :name,
			"price" = :price,
			"barcode" = :barcode,
			"category_id" = :category_id,
			"updated_at" = NOW()
	`

	params = map[string]interface{}{
		"id":          req.Id,
		"name":        req.Name,
		"price":       req.Price,
		"barcode":     req.Barcode,
		"category_id": req.CategoryId,
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
		return "", fmt.Errorf("product with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (r *productRepo) Delete(req *models.ProductPrimaryKey) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx, "DELETE FROM product WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product with ID %s not found", req.Id)

	}

	return nil
}

// get by barcode
func (r *productRepo) GetByBarcode(req *models.ProductBarcodeRequest) (*models.ProductBarcodeResponse, error) {

	var (
		name        sql.NullString
		price       sql.NullFloat64
		category_id sql.NullString
	)

	query := `
		SELECT
			"name",
			"price",		
			"category_id"
		FROM "product"
		WHERE "barcode" = $1
	`

	err := r.db.QueryRow(context.Background(), query, req.Barcode).Scan(
		&name,
		&price,
		&category_id,
	)

	if err != nil {
		return nil, err
	}

	return &models.ProductBarcodeResponse{
		Name:       name.String,
		Price:      price.Float64,
		CategoryId: category_id.String,
	}, nil
}
