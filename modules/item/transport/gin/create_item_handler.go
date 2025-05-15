package ginItem

import (
	"golang_todo_list_api/common"
	"golang_todo_list_api/modules/item/business"
	"golang_todo_list_api/modules/item/model"
	"golang_todo_list_api/modules/item/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateItems(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var data model.TodoItemCreation
		if err := ctx.ShouldBind(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := storage.NewSQLStore(db)

		biz := business.NewCreateItemBussines(store)
		if err := biz.CreateNewItem(ctx.Request.Context(), &data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
