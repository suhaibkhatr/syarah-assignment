package main

import (
	"gift-store/internal/config"
	"gift-store/internal/db"
	"gift-store/internal/repo/mysql"
	"gift-store/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()
	source := db.NewSource(config)
	productRepo := mysql.NewProductRepo(source)
	productService := services.NewProductService(productRepo)
	g := gin.Default()
	v1 := g.Group("/api/v1")
	{
		v1.GET("/products", productService.GetAllProducts)
		v1.GET("/products/:id", productService.GetProductByID)
		v1.POST("/products", productService.Insert)
		v1.PUT("/products/:id", productService.Update)
		v1.DELETE("/products/:id", productService.Delete)
	}
	g.Run(":" + config.APP_PORT)
}
