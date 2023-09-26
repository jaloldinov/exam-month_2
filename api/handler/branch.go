package handler

import (
	"market/models"
	"market/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBranch godoc
// @Router       /branch [POST]
// @Summary      CREATE BRANCH
// @Description adds branch data to db based on given info in body
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateBranch  true  "branch data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) CreateBranch(ctx *gin.Context) {
	var branch models.CreateBranch
	err := ctx.ShouldBind(&branch)
	if err != nil {
		h.log.Error("error while binding branch:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.strg.Branch().Create(&branch)
	if err != nil {
		h.log.Error("error branch create:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "success", "resp": resp})
}

// ListBranchs godoc
// @Router       /branch [GET]
// @Summary      LIST BRANCHS
// @Description  gets all branch based on limit, page and search by name
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 search        query     string     false  "search"
// @Success      200  {object}  models.BranchGetListResponse
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetListBranch(ctx *gin.Context) {
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

	resp, err := h.strg.Branch().GetList(&models.BranchGetListRequest{
		Page:   page,
		Limit:  limit,
		Search: ctx.Query("search"),
	})
	if err != nil {
		h.log.Error("error Branch GetListBranch:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetBranch godoc
// @Router       /branch/{id} [GET]
// @Summary      GET BY ID
// @Description  gets branch by ID
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Branch ID" format(uuid)
// @Success      200  {object}  models.Branch
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetByIDBranch(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.strg.Branch().GetByID(&models.BranchPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error get branch:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateBranch godoc
// @Router       /branch/{id} [PUT]
// @Summary      UPDATE BRANCH
// @Description  UPDATES BRANCH BASED ON GIVEN DATA AND ID
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of branch" format(uuid)
// @Param        data  body      models.CreateBranch  true  "branch data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) UpdateBranch(ctx *gin.Context) {
	var branch models.UpdateBranch

	err := ctx.ShouldBind(&branch)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	branch.Id = ctx.Param("id")
	resp, err := h.strg.Branch().Update(&branch)
	if err != nil {
		h.log.Error("error branch update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteBranch godoc
// @Router       /branch/{id} [DELETE]
// @Summary      DELETE BRANCH BY ID
// @Description  deletes branch by id
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of branch" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) DeleteBranch(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.strg.Branch().Delete(&models.BranchPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error deleting branch:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}
