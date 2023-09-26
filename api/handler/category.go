package handler

import (
	"market/models"
	"market/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Router       /category [POST]
// @Summary      CREATE CATEGORY
// @Description adds category data to db based on given info in body
// @Tags         CATEGORY
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateCategory  true  "category data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) CreateCategory(ctx *gin.Context) {
	var category models.CreateCategory
	err := ctx.ShouldBind(&category)
	if err != nil {
		h.log.Error("error while binding category:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.strg.Category().Create(&category)
	if err != nil {
		h.log.Error("error category create:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "success", "resp": resp})
}

// ListCategorys godoc
// @Router       /category [GET]
// @Summary      LIST CATEGORY
// @Description  gets all category based on limit, page and search by name
// @Tags         CATEGORY
// @Accept       json
// @Produce      json
// @Param  		 limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param  		 page          query     int        false  "page"           minimum(1)     default(1)
// @Param   	 search        query     string     false  "search"
// @Success      200  {object}  models.CategoryGetListResponse
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetListCategory(ctx *gin.Context) {
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

	resp, err := h.strg.Category().GetList(&models.CategoryGetListRequest{
		Page:   page,
		Limit:  limit,
		Search: ctx.Query("search"),
	})
	if err != nil {
		h.log.Error("error Category GetListCategory:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetCategory godoc
// @Router       /category/{id} [GET]
// @Summary      GET BY ID
// @Description  gets category by ID
// @Tags         CATEGORY
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID" format(uuid)
// @Success      200  {object}  models.Category
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) GetByIDCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.strg.Category().GetByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error get category:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateCategory godoc
// @Router       /category/{id} [PUT]
// @Summary      UPDATE CATEGORY
// @Description  UPDATES CATEGORY BASED ON GIVEN DATA AND ID
// @Tags         CATEGORY
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of category" format(uuid)
// @Param        data  body      models.CreateCategory  true  "category data"
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	var category models.UpdateCategory

	err := ctx.ShouldBind(&category)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	category.Id = ctx.Param("id")
	resp, err := h.strg.Category().Update(&category)
	if err != nil {
		h.log.Error("error category update:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "resp": resp})
}

// DeleteCategory godoc
// @Router       /category/{id} [DELETE]
// @Summary      DELETE CATEGORY BY ID
// @Description  deletes category by id
// @Tags         CATEGORY
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of category" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  models.ErrorResp
// @Failure      404  {object}  models.ErrorResp
// @Failure      500  {object}  models.ErrorResp
func (h *Handler) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.strg.Category().Delete(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.log.Error("error deleting category:", logger.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success"})
}
