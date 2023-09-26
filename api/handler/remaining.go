package handler

import (
	"market/models"
	"market/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRemaining godoc
// @Router       /do_income/{coming_table_id} [POST]
// @Summary      CREATE REMAINING
// @Description adds remaining data to db based on given coming_table_id
// @Tags         REMAINING
// @Accept       json
// @Produce      json
// @Param        coming_table_id path string true "Coming Table ID"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) CreateRemaining(ctx *gin.Context) {

	comingTableID := ctx.Param("coming_table_id")

	var remaining models.CreateRemaining

	// checking coming_table info weather it is in_process or finished // return branch_id
	comingTableId := models.ComingTablePrimaryKey{Id: comingTableID}
	branch_id, err := h.strg.ComingTable().GetStatus(&comingTableId)
	if err != nil {
		h.log.Error("error while getting coming table status", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	remaining.BranchId = branch_id

	// get coming_table_product details using coming_id
	coming_id := models.ComingTableProductPrimaryKey{Id: comingTableID}
	coming_table_data, err := h.strg.ComingTableProduct().GetByComingTableId(&coming_id)
	if err != nil {
		h.log.Error("error while getting coming_table_data details:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not Found Coming_table_data with that coming_id"})
		return
	}
	remaining.Barcode = coming_table_data.ProductBarcode
	remaining.CategoryId = coming_table_data.CategoryId
	remaining.Price = coming_table_data.ProductPrice
	remaining.Count = coming_table_data.Count
	remaining.TotalPrice = coming_table_data.TotalPrice
	remaining.Name = coming_table_data.ProductName

	// cheking remain table by branch and barcode, if yes, update otherwise create
	req := models.CheckingRemaining{BranchId: branch_id, Barcode: remaining.Barcode}

	id, err := h.strg.Remaining().CheckRemaing(&req)
	if err != nil {
		h.log.Info("not found then we are creating new remaining", logger.Error(err))
		// if product or remaining is not exists, add remaining table
		resp, err := h.strg.Remaining().Create(&remaining)
		if err != nil {
			h.log.Error("error remaining create:", logger.Error(err))
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "added new remaining", "resp": resp})
		// change status
		h.strg.ComingTable().UpdateStatus(&comingTableId)
		return
	}

	// if exits Update Coming Product Table
	var updatingData = models.UpdateRemaining{
		Id:         id,
		BranchId:   remaining.BranchId,
		CategoryId: remaining.CategoryId,
		Name:       remaining.Name,
		Price:      remaining.Price,
		Barcode:    remaining.Barcode,
		Count:      remaining.Count,
		TotalPrice: remaining.TotalPrice,
	}
	r, err := h.strg.Remaining().UpdateExists(&updatingData)
	if err != nil {
		h.log.Info("error remaing update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "updated existing remaining table", "resp": r})

	// if everything ok, change status to finish
	h.strg.ComingTable().UpdateStatus(&comingTableId)
}

// ListRemainings godoc
// @Router       /remaining [GET]
// @Summary      LIST REMAINING
// @Description  gets all remaining based on limit, page and search by name
// @Tags         REMAINING
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 branch_id          query     string     false  "branch_id"
// @Param   	 category_id        query     string     false  "category_id"
// @Param   	 barcode            query     string     false  "barcode"
// @Success      200  {object}  models.RemainingGetListResponse
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetListRemaining(ctx *gin.Context) {
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

	resp, err := h.strg.Remaining().GetList(&models.RemainingGetListRequest{
		Page:       page,
		Limit:      limit,
		CategoryId: ctx.Query("search"),
		Barcode:    ctx.Query("barcode"),
		BranchId:   ctx.Query("branch_id"),
	})
	if err != nil {
		h.log.Error("error Remaining GetListRemaining:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetRemaining godoc
// @Router       /remaining/{id} [GET]
// @Summary      GET BY ID
// @Description  gets remaining by ID
// @Tags         REMAINING
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Remaining ID" format(uuid)
// @Success      200  {object}  models.Remaining
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetByIDRemaining(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.strg.Remaining().GetByID(&models.RemainingPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not Found Product"})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateRemaining godoc
// @Router       /remaining/{id} [PUT]
// @Summary      UPDATE REMAINING
// @Description  UPDATES REMAINING BASED ON GIVEN DATA AND ID
// @Tags         REMAINING
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of remaining" format(uuid)
// @Param        data  body      models.UpdateRemainingSoft  true  "remaining data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) UpdateRemaining(ctx *gin.Context) {
	var remaining models.UpdateRemaining

	err := ctx.ShouldBind(&remaining)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	remaining.Id = ctx.Param("id")
	remaining.TotalPrice = float64(remaining.Count) * remaining.Price
	resp, err := h.strg.Remaining().Update(&remaining)
	if err != nil {
		h.log.Error("error remaining update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteRemaining godoc
// @Router       /remaining/{id} [DELETE]
// @Summary      DELETE REMAINING BY ID
// @Description  deletes remaining by id
// @Tags         REMAINING
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of remaining" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) DeleteRemaining(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.strg.Remaining().Delete(&models.RemainingPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error deleting remaining:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}
