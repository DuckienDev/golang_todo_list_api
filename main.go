package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TodoItem struct {
	Id          int        `json:"id" gorm:"column:id"`
	Title       string     `json:"title" gorm:"column:title"`
	Description string     `json:"description" gorm:"column:description"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (TodoItem) TableName() string { return "todo_items" }

type TodoItemCreation struct {
	Id          int    `json:"-" gorm:"column:id"`
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
	Status      string `json:"status" gorm:"column:status"`
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title"`
	Description *string `json:"description" gorm:"column:description"`
	Status      *string `json:"status" gorm:"column:status"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"`
}

func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit >= 100 {
		p.Limit = 10
	}
}

func main() {

	godotenv.Load()

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	// CRUD : Create , Read , Update, Delete\
	// POST /v1/items/  (create new item)
	// GET /v1/items/   (list items)  /v1/items?page=1
	// GET /v1/items/:id (get item detail by id)
	// (PUT  || PATCH) /v1/items/:id  (update an item by id)
	// DELETE /v1/items/:id  (delete item by id)

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItems(db))
			items.GET("", ListItem(db))
			items.GET("/:id", GetItem(db))
			items.PATCH("/:id", UpdateItem(db))
			items.DELETE("/:id", DeleteItem(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Connected success",
		})
	})
	r.Run(":3000")
}

func CreateItems(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var data TodoItemCreation
		if err := ctx.ShouldBind(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if data.Status == "" {
			data.Status = "Doing"
		}

		if err := db.Create(&data).Error; err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": data.Id,
		})
	}
}

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItem
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemUpdate
		id, err := strconv.Atoi(c.Param("id"))
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
			"status": "Delete",
		}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {

		var paging Paging

		if err := ctx.ShouldBind(&paging); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		paging.Process()

		var result []TodoItem

		db = db.Where("status <> ?", "Delete")

		if err := db.Table(TodoItem{}.TableName()).
			Count(&paging.Total).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Order("id desc").
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&result).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":   result,
			"paging": paging,
		})
	}
}
