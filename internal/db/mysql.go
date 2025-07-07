package db

import (
	"database/sql"
	"gift-store/internal/config"
	"gift-store/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLSource struct {
	db *sql.DB
}

func NewSource(cfg config.AppConfig) *MySQLSource {
	db, err := sql.Open("mysql", cfg.MySQL.DSN)
	if err != nil {
		panic(err)
	}
	return &MySQLSource{db: db}
}

func (s *MySQLSource) GetAllProducts() ([]models.Product, error) {
	rows, err := s.db.Query("SELECT id, name, description, price, category, gift_category, age_group, brand, is_available, created_at, updated_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Category, &p.GiftCategory, &p.AgeGroup, &p.Brand, &p.IsAvailable, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			continue
		}
		products = append(products, p)
	}
	return products, nil
}
