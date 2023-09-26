package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"market/models"
	"market/pkg/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type comingTableProduct struct {
	db *pgxpool.Pool
}

func NewComingTableProductRepo(db *pgxpool.Pool) *comingTableProduct {
	return &comingTableProduct{
		db: db,
	}
}

func (r *comingTableProduct) Create(req *models.CreateComingTableProduct) (string, error) {
	var (
		id = uuid.NewString()
	)

	query := `
				INSERT INTO "coming_table_product"(
					"id",
					"category_id",
					"name",
					"price",
					"barcode",
					"count",
					"total_price",
					"coming_table_id",
					"created_at")
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8,  NOW())`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.CategoryId,
		req.ProductName,
		req.ProductPrice,
		req.ProductBarcode,
		req.Count,
		req.TotalPrice,
		req.ComingTableId,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *comingTableProduct) GetByID(req *models.ComingTableProductPrimaryKey) (*models.ComingTableProduct, error) {

	var (
		id              sql.NullString
		category_id     sql.NullString
		name            sql.NullString
		price           sql.NullFloat64
		barcode         sql.NullString
		count           sql.NullInt16
		total_price     sql.NullFloat64
		coming_table_id sql.NullString
		created_at      sql.NullString
		updated_at      sql.NullString
	)

	query := `
			SELECT
					"id",
					"category_id",
					"name",
					"price",
					"barcode",
					"count",
					"total_price",
					"coming_table_id",
					"created_at"
					"updated_at" 
			FROM "coming_table_product"
			WHERE id = $1 `

	err := r.db.QueryRow(context.Background(), query, req.Id).Scan(
		&id,
		&category_id,
		&name,
		&price,
		&barcode,
		&count,
		&total_price,
		&coming_table_id,
		&created_at,
		&updated_at,
	)

	if err != nil {
		return nil, err
	}

	return &models.ComingTableProduct{
		Id:             id.String,
		CategoryId:     category_id.String,
		ProductName:    name.String,
		ProductPrice:   price.Float64,
		ProductBarcode: barcode.String,
		Count:          int(count.Int16),
		TotalPrice:     total_price.Float64,
		ComingTableId:  coming_table_id.String,
		CreatedAt:      created_at.String,
		UpdatedAt:      updated_at.String,
	}, nil
}

func (r *comingTableProduct) GetList(req *models.ComingTableProductGetListRequest) (*models.ComingTableProductGetListResponse, error) {
	params := make(map[string]interface{})
	var resp = &models.ComingTableProductGetListResponse{}

	resp.ComingTableProducts = make([]*models.ComingTableProduct, 0)

	filter := " WHERE true "
	query := `
			SELECT
				COUNT(*) OVER(),
				"id",
				"category_id",
				"name",
				"price",
				"barcode",
				"count",
				"total_price",
				"coming_table_id",
				"created_at",
				"updated_at" 
			FROM "coming_table_product"
		`
	if req.CategoryId != "" {
		filter += ` AND ("category_id" = :category_id)`
		params["category_id"] = req.CategoryId
	}

	if req.ProductBarcode != "" {
		filter += ` AND ("barcode" = :barcode)`
		params["barcode"] = req.ProductBarcode
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
			id              sql.NullString
			category_id     sql.NullString
			name            sql.NullString
			price           sql.NullFloat64
			barcode         sql.NullString
			count           sql.NullInt16
			total_price     sql.NullFloat64
			coming_table_id sql.NullString
			created_at      sql.NullString
			updated_at      sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&category_id,
			&name,
			&price,
			&barcode,
			&count,
			&total_price,
			&coming_table_id,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		resp.ComingTableProducts = append(resp.ComingTableProducts, &models.ComingTableProduct{
			Id:             id.String,
			CategoryId:     category_id.String,
			ProductName:    name.String,
			ProductPrice:   price.Float64,
			ProductBarcode: barcode.String,
			Count:          int(count.Int16),
			TotalPrice:     total_price.Float64,
			ComingTableId:  coming_table_id.String,
			CreatedAt:      created_at.String,
			UpdatedAt:      updated_at.String,
		})
	}
	return resp, nil
}

func (r *comingTableProduct) Update(req *models.UpdateComingTableProduct) (string, error) {

	query := `
		UPDATE
			"coming_table_product"
		SET
				"category_id" = $1,
				"name" = $2,
				"price" = $3,
				"barcode" = $4,
				"count" = $5,
				"total_price" = $6,
				"coming_table_id" = $7,
				"updated_at" = NOW()
				WHERE id = $8
	`

	result, err := r.db.Exec(context.Background(), query,
		req.CategoryId,
		req.ProductName,
		req.ProductPrice,
		req.ProductBarcode,
		req.Count,
		req.TotalPrice,
		req.ComingTableId,
		req.Id,
	)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("coming_table_product with ID %s not found", req.Id)
	}

	return req.Id, nil
}

func (r *comingTableProduct) Delete(req *models.ComingTableProductPrimaryKey) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx, "DELETE FROM coming_table_product WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("coming_table_product with ID %s not found", req.Id)

	}

	return nil
}

func (r *comingTableProduct) CheckExistProduct(req *models.ComingTableProductBarcode) (string, error) {
	var id sql.NullString

	query := `
		SELECT
			"id"
		FROM "coming_table_product"
		WHERE "barcode" = $1 and "coming_table_id" = $2`

	err := r.db.QueryRow(context.Background(), query, req.Barcode, req.ComingTableId).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("not found")
		}
		return "", err
	}

	return id.String, nil
}

func (r *comingTableProduct) UpdateIdExists(req *models.UpdateComingTableProduct) (string, error) {
	query := `
		UPDATE
			"coming_table_product"
		SET
			"category_id" = $1,
			"barcode" = $2,
			"name" = $3,
			"price" = $4,
			"count" = "count" + $5,
			"total_price" = "total_price" + $6,
			"coming_table_id" = $7,
			"updated_at" = NOW()
		WHERE
			"id" = $8
	`

	result, err := r.db.Exec(context.Background(), query,
		req.CategoryId,
		req.ProductBarcode,
		req.ProductName,
		req.ProductPrice,
		req.Count,
		req.TotalPrice,
		req.ComingTableId,
		req.Id,
	)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("coming_table_product with id %s not found", req.ComingTableId)
	}

	return req.ComingTableId, nil
}

func (r *comingTableProduct) GetByComingTableId(req *models.ComingTableProductPrimaryKey) (*models.ComingTableProduct, error) {

	var (
		id          sql.NullString
		category_id sql.NullString
		name        sql.NullString
		price       sql.NullFloat64
		barcode     sql.NullString
		count       sql.NullInt16
		total_price sql.NullFloat64
	)

	query := `
			SELECT
					"id",
					"category_id",
					"name",
					"price",
					"barcode",
					sum("count"),
					sum("total_price")
			FROM "coming_table_product"
			WHERE "coming_table_id" = $1 
			GROUP BY "id", "barcode"
			`

	err := r.db.QueryRow(context.Background(), query, req.Id).Scan(
		&id,
		&category_id,
		&name,
		&price,
		&barcode,
		&count,
		&total_price,
	)

	if err != nil {
		return nil, err
	}

	return &models.ComingTableProduct{
		Id:             id.String,
		CategoryId:     category_id.String,
		ProductName:    name.String,
		ProductPrice:   price.Float64,
		ProductBarcode: barcode.String,
		Count:          int(count.Int16),
		TotalPrice:     total_price.Float64,
	}, nil
}
