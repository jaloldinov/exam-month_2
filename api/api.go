package api

import (
	_ "market/api/docs"
	"market/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("/branch", h.CreateBranch)
	r.GET("/branch/:id", h.GetByIDBranch)
	r.GET("/branch", h.GetListBranch)
	r.PUT("/branch/:id", h.UpdateBranch)
	r.DELETE("/branch/:id", h.DeleteBranch)

	r.POST("/category", h.CreateCategory)
	r.GET("/category/:id", h.GetByIDCategory)
	r.GET("/category", h.GetListCategory)
	r.PUT("/category/:id", h.UpdateCategory)
	r.DELETE("/category/:id", h.DeleteCategory)

	r.POST("/product", h.CreateProduct)
	r.GET("/product/:id", h.GetByIDProduct)
	r.GET("/product", h.GetListProduct)
	r.PUT("/product/:id", h.UpdateProduct)
	r.DELETE("/product/:id", h.DeleteProduct)

	r.POST("/coming_table", h.CreateComingTable)
	r.GET("/coming_table/:id", h.GetByIDComingTable)
	r.GET("/coming_table", h.GetListComingTable)
	r.PUT("/coming_table/:id", h.UpdateComingTable)
	r.DELETE("/coming_table/:id", h.DeleteComingTable)

	r.POST("/coming_product/:coming_table_id", h.CreateComingTableProduct)

	r.GET("/coming_product/:id", h.GetByIDComingTableProduct)
	r.GET("/coming_product", h.GetListComingTableProduct)
	r.PUT("/coming_product/:id", h.UpdateComingTableProduct)
	r.DELETE("/coming_product/:id", h.DeleteComingTableProduct)

	r.POST("/do_income/:coming_table_id", h.CreateRemaining)

	r.GET("/remaining/:id", h.GetByIDRemaining)
	r.GET("/remaining", h.GetListRemaining)
	r.PUT("/remaining/:id", h.UpdateRemaining)
	r.DELETE("/remaining/:id", h.DeleteRemaining)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
