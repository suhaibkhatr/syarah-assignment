package mysql

import (
	"gift-store/internal/db"
	"gift-store/internal/models"
)

type ProductRepo struct {
	source *db.MySQLSource
}

func NewProductRepo(source *db.MySQLSource) *ProductRepo {
	return &ProductRepo{source: source}
}

func (s *ProductRepo) GetAllProducts() ([]models.Product, error) {
	rows, err := s.source.Db.Query("SELECT id, name, description, price, category, gift_category, age_group, brand, is_available, created_at, updated_at FROM products")
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
