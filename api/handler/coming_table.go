package handler

import (
	"market/models"
	"market/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComingTable godoc
// @Router       /coming_table [POST]
// @Summary      CREATE COMING TABLE
// @Description adds coming_table data to db based on given info in body
// @Tags         COMING TABLE
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateComingTable  true  "coming_table data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) CreateComingTable(ctx *gin.Context) {
	var coming_table models.CreateComingTable
	err := ctx.ShouldBind(&coming_table)
	if err != nil {
		h.log.Error("error while binding coming_table:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.strg.ComingTable().Create(&coming_table)
	if err != nil {
		h.log.Error("error coming_table create:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "success", "resp": resp})
}

// ListComingTables godoc
// @Router       /coming_table [GET]
// @Summary      LIST COMING TABLES
// @Description  gets all coming_table based on limit, page and search by name
// @Tags         COMING TABLE
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 coming_id        query     string     false  "coming_id"
// @Param   	 branch_id        query     string     false  "branch_id"
// @Success      200  {object}  models.ComingTableGetListResponse
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetListComingTable(ctx *gin.Context) {
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

	resp, err := h.strg.ComingTable().GetList(&models.ComingTableGetListRequest{
		Page:     page,
		Limit:    limit,
		ComingId: ctx.Query("coming_id"),
		BranchId: ctx.Query("branch_id"),
	})
	if err != nil {
		h.log.Error("error ComingTable GetListComingTable:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetComingTable godoc
// @Router       /coming_table/{id} [GET]
// @Summary      GET BY ID
// @Description  gets coming_table by ID
// @Tags         COMING TABLE
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ComingTable ID" format(uuid)
// @Success      200  {object}  models.ComingTable
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetByIDComingTable(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.strg.ComingTable().GetByID(&models.ComingTablePrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error get coming_table:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateComingTable godoc
// @Router       /coming_table/{id} [PUT]
// @Summary      UPDATE COMING TABLE
// @Description  UPDATES COMING TABLE BASED ON GIVEN DATA AND ID
// @Tags         COMING TABLE
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of coming_table" format(uuid)
// @Param        data  body      models.CreateComingTable  true  "coming_table data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) UpdateComingTable(ctx *gin.Context) {
	var coming_table models.UpdateComingTable

	err := ctx.ShouldBind(&coming_table)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	coming_table.Id = ctx.Param("id")
	resp, err := h.strg.ComingTable().Update(&coming_table)
	if err != nil {
		h.log.Error("error coming_table update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteComingTable godoc
// @Router       /coming_table/{id} [DELETE]
// @Summary      DELETE COMING TABLE BY ID
// @Description  deletes coming_table by id
// @Tags         COMING TABLE
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of coming_table" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) DeleteComingTable(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.strg.ComingTable().Delete(&models.ComingTablePrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error deleting coming_table:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}
