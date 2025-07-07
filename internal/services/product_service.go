package services

import (
	"gift-store/internal/models"
	"gift-store/internal/repo/mysql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductService struct {
	productRepo *mysql.ProductRepo
}

func NewProductService(productRepo *mysql.ProductRepo) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetAllProducts(c *gin.Context) {
	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (s *ProductService) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}
	product, err := s.productRepo.GetProductByID(idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, product)
}

func (s *ProductService) Insert(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id, err := s.productRepo.Insert(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	p.ID = id
	c.JSON(http.StatusOK, gin.H{"message": "Product inserted successfully", "product": p})
}

func (s *ProductService) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = s.productRepo.Update(idInt, &p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	p.ID = idInt
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": p})
}

func (s *ProductService) Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}
	err = s.productRepo.Delete(idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
