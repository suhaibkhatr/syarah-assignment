package mysql

import (
	"fmt"
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

func (s *ProductRepo) GetProductByID(id int) (*models.Product, error) {
	row := s.source.Db.QueryRow("SELECT id, name, description, price, category, gift_category, age_group, brand, is_available, created_at, updated_at FROM products WHERE id = ?", id)
	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Category, &p.GiftCategory, &p.AgeGroup, &p.Brand, &p.IsAvailable, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ProductRepo) Insert(p *models.Product) (int, error) {
	// For MySQL, we use LastInsertId() to get the auto-incremented ID
	result, err := s.source.Db.Exec(
		"INSERT INTO products (name, description, price, category, gift_category, age_group, brand, is_available, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		p.Name, p.Description, p.Price, p.Category, p.GiftCategory, p.AgeGroup, p.Brand, p.IsAvailable, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert product: %w", err)
	}

	// Get the ID of the newly inserted product
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(id), nil
}

func (s *ProductRepo) Update(id int, p *models.Product) error {
	_, err := s.source.Db.Exec("UPDATE products SET name = ?, description = ?, price = ?, category = ?, gift_category = ?, age_group = ?, brand = ?, is_available = ?, created_at = ?, updated_at = ? WHERE id = ?", p.Name, p.Description, p.Price, p.Category, p.GiftCategory, p.AgeGroup, p.Brand, p.IsAvailable, p.CreatedAt, p.UpdatedAt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductRepo) Delete(id int) error {
	_, err := s.source.Db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
