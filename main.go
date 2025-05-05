package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoItem struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Stattus     string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func main() {
	fmt.Println("Hello, World!")

	now := time.Now().UTC()

	item := TodoItem{
		ID:          1,
		Title:       "This is item 1",
		Description: "Decription 1",
		Stattus:     "Doing",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": item,
		})
	})
	r.Run(":3000")
}
