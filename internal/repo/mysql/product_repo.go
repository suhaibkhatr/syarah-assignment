package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"gift-store/internal/db"
	"gift-store/internal/models"
	"time"
)

type ProductRepo struct {
	source *db.MySQLSource
}

func NewProductRepo(source *db.MySQLSource) *ProductRepo {
	return &ProductRepo{source: source}
}

func (s *ProductRepo) GetAllProducts() ([]models.Product, error) {
	fmt.Println("GetAllProducts")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT id, name, description, price, category, gift_category, age_group, brand, is_available, created_at, updated_at FROM products`
	rows, err := s.source.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Category,
			&p.GiftCategory,
			&p.AgeGroup,
			&p.Brand,
			&p.IsAvailable,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			fmt.Printf("Error scanning product row: %v\n", err)
			continue
		}

		// Parse the time strings
		layout := "2006-01-02 15:04:05" // MySQL datetime format
		if createdAtStr != "" {
			createdAt, err := time.Parse(layout, createdAtStr)
			if err != nil {
				fmt.Printf("Error parsing created_at: %v\n", err)
			} else {
				p.CreatedAt = createdAt
			}
		}

		if updatedAtStr != "" {
			updatedAt, err := time.Parse(layout, updatedAtStr)
			if err != nil {
				fmt.Printf("Error parsing updated_at: %v\n", err)
			} else {
				p.UpdatedAt = updatedAt
			}
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}

func (s *ProductRepo) GetProductByID(id int) (*models.Product, error) {
	query := "SELECT id, name, description, price, category, gift_category, age_group, brand, is_available, created_at, updated_at FROM products WHERE id = ?"
	row := s.source.Db.QueryRow(query, id)

	var p models.Product
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Category,
		&p.GiftCategory,
		&p.AgeGroup,
		&p.Brand,
		&p.IsAvailable,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("error scanning product: %w", err)
	}

	// Parse the time strings
	layout := "2006-01-02 15:04:05" // MySQL datetime format
	if createdAtStr != "" {
		p.CreatedAt, err = time.Parse(layout, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %w", err)
		}
	}

	if updatedAtStr != "" {
		p.UpdatedAt, err = time.Parse(layout, updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing updated_at: %w", err)
		}
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
