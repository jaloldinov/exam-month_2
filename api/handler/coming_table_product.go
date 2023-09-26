package handler

import (
	"market/models"
	"market/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComingTableProduct godoc
// @Router       /coming_product/{coming_table_id} [POST]
// @Summary      CREATE COMING TABLE PRODUCT
// @Description adds coming_product data to db based on given info in body
// @Tags         COMING TABLE PRODUCT
// @Accept       json
// @Produce      json
// @Param        coming_table_id path string true "Coming Table ID"
// @Param        barcode query string true "Barcode value"
// @Param        data  body      models.CreateComingTableProductCount  true  "coming_product count"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) CreateComingTableProduct(ctx *gin.Context) {

	comingTableID := ctx.Param("coming_table_id")
	barcodeQ := ctx.Query("barcode")

	var coming_product models.CreateComingTableProduct
	err := ctx.ShouldBind(&coming_product)
	if err != nil {
		h.log.Error("error while binding coming_product:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	// checking coming_table info weather it is in_process or finished
	comingTableId := models.ComingTablePrimaryKey{Id: comingTableID}
	_, err = h.strg.ComingTable().GetStatus(&comingTableId)
	if err != nil {
		h.log.Error("error while getting coming table status", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// get product details (name, price, category_id)
	productBarcode := models.ProductBarcodeRequest{Barcode: barcodeQ}
	productDetails, err := h.strg.Product().GetByBarcode(&productBarcode)
	if err != nil {
		h.log.Error("error while getting product details:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not Found Product with that barcode"})
		return
	}

	// filling all requesting data to coming table
	coming_product.CategoryId = productDetails.CategoryId
	coming_product.ProductName = productDetails.Name
	coming_product.ProductPrice = productDetails.Price
	coming_product.ProductBarcode = barcodeQ
	coming_product.TotalPrice = (productDetails.Price * float64(coming_product.Count))
	coming_product.ComingTableId = comingTableID

	//  Checking exists product by shtrixcode in coming_table_product table
	barcode := models.ComingTableProductBarcode{Barcode: barcodeQ, ComingTableId: comingTableID}
	id, err := h.strg.ComingTableProduct().CheckExistProduct(&barcode)
	if err != nil {
		h.log.Info("product or coming_table_id not found:", logger.Error(err))

		// if product or coming_table_id is not exists, ADD Coming product table
		resp, err := h.strg.ComingTableProduct().Create(&coming_product)
		if err != nil {
			h.log.Error("error coming_product create:", logger.Error(err))
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "added coming_product_table", "resp": resp})

		return
	}

	// if exits Update Coming Product Table
	var updatingData = models.UpdateComingTableProduct{
		Id:             id,
		CategoryId:     coming_product.CategoryId,
		ProductName:    coming_product.ProductName,
		ProductPrice:   coming_product.ProductPrice,
		ProductBarcode: barcodeQ,
		Count:          coming_product.Count,
		TotalPrice:     (productDetails.Price * float64(coming_product.Count)),
		ComingTableId:  comingTableID,
	}
	r, err := h.strg.ComingTableProduct().UpdateIdExists(&updatingData)
	if err != nil {
		h.log.Info("error coming_table_product update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "updated existing coming_product_table", "resp": r})
}

// ListComingTableProducts godoc
// @Router       /coming_product [GET]
// @Summary      LIST COMING TABLE PRODUCT
// @Description  gets all coming_product based on limit, page and search by name
// @Tags         COMING TABLE PRODUCT
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 category_id        query     string     false  "category_id"
// @Param   	 barcode            query     string     false  "barcode"
// @Success      200  {object}  models.ComingTableProductGetListResponse
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetListComingTableProduct(ctx *gin.Context) {
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

	resp, err := h.strg.ComingTableProduct().GetList(&models.ComingTableProductGetListRequest{
		Page:           page,
		Limit:          limit,
		CategoryId:     ctx.Query("category_id"),
		ProductBarcode: ctx.Query("barcode"),
	})
	if err != nil {
		h.log.Error("error ComingTableProduct GetListComingTableProduct:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetComingTableProduct godoc
// @Router       /coming_product/{id} [GET]
// @Summary      GET BY ID
// @Description  gets coming_product by ID
// @Tags         COMING TABLE PRODUCT
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ComingTableProduct ID" format(uuid)
// @Success      200  {object}  models.ComingTableProduct
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetByIDComingTableProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.strg.ComingTableProduct().GetByID(&models.ComingTableProductPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not Found Product"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateComingTableProduct godoc
// @Router       /coming_product/{id} [PUT]
// @Summary      UPDATE COMING TABLE PRODUCT
// @Description  UPDATES COMING TABLE PRODUCT BASED ON GIVEN DATA AND ID
// @Tags         COMING TABLE PRODUCT
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of coming_product" format(uuid)
// @Param        data  body      models.CreateComingTableProduct  true  "coming_product data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) UpdateComingTableProduct(ctx *gin.Context) {
	var coming_product models.UpdateComingTableProduct

	err := ctx.ShouldBind(&coming_product)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	coming_product.Id = ctx.Param("id")
	resp, err := h.strg.ComingTableProduct().Update(&coming_product)
	if err != nil {
		h.log.Error("error coming_product update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteComingTableProduct godoc
// @Router       /coming_product/{id} [DELETE]
// @Summary      DELETE COMING TABLE PRODUCT BY ID
// @Description  deletes coming_product by id
// @Tags         COMING TABLE PRODUCT
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of coming_product" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) DeleteComingTableProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.strg.ComingTableProduct().Delete(&models.ComingTableProductPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error deleting coming_product:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}
