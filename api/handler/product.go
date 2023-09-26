package handler

import (
	"market/models"
	"market/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Router       /product [POST]
// @Summary      CREATE PRODUCT
// @Description adds product data to db based on given info in body
// @Tags         PRODUCT
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateProduct  true  "product data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) CreateProduct(ctx *gin.Context) {
	var product models.CreateProduct
	err := ctx.ShouldBind(&product)
	if err != nil {
		h.log.Error("error while binding product:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.strg.Product().Create(&product)
	if err != nil {
		h.log.Error("error product create:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "success", "resp": resp})
}

// ListProducts godoc
// @Router       /product [GET]
// @Summary      LIST PRODUCT
// @Description  gets all product based on limit, page and search by name
// @Tags         PRODUCT
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 barcode        query     string     false  "barcode"
// @Param   	 name        query     string     false  "name"
// @Success      200  {object}  models.ProductGetListResponse
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetListProduct(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		h.log.Error("error get page:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error get limit:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid page param")
		return
	}

	resp, err := h.strg.Product().GetList(&models.ProductGetListRequest{
		Page:    page,
		Limit:   limit,
		Name:    ctx.Query("name"),
		Barcode: ctx.Query("barcode"),
	})
	if err != nil {
		h.log.Error("error Product GetListProduct:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetProduct godoc
// @Router       /product/{id} [GET]
// @Summary      GET BY ID
// @Description  gets product by ID
// @Tags         PRODUCT
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID" format(uuid)
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetByIDProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.strg.Product().GetByID(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error get product:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not Found Product"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateProduct godoc
// @Router       /product/{id} [PUT]
// @Summary      UPDATE PRODUCT
// @Description  UPDATES PRODUCT BASED ON GIVEN DATA AND ID
// @Tags         PRODUCT
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of product" format(uuid)
// @Param        data  body      models.CreateProduct  true  "product data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var product models.UpdateProduct

	err := ctx.ShouldBind(&product)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	product.Id = ctx.Param("id")
	resp, err := h.strg.Product().Update(&product)
	if err != nil {
		h.log.Error("error product update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteProduct godoc
// @Router       /product/{id} [DELETE]
// @Summary      DELETE PRODUCT BY ID
// @Description  deletes product by id
// @Tags         PRODUCT
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of product" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.strg.Product().Delete(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error deleting product:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}
