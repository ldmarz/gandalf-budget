package store

import (
	"database/sql" // For sql.ErrNoRows
	"fmt"
	"log"
)

func (s *sqlStore) GetAllCategories() ([]Category, error) {
	var categories []Category
	err := s.DB.Select(&categories, "SELECT id, name, color FROM categories ORDER BY name ASC")
	if err != nil {
		log.Printf("Error getting all categories: %v", err)
		return nil, err
	}
	if categories == nil {
		categories = []Category{}
	}
	return categories, nil
}

func (s *sqlStore) CreateCategory(category *Category) error {
	if category.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	if category.Color == "" {
		return fmt.Errorf("category color cannot be empty")
	}
	query := `INSERT INTO categories (name, color) VALUES (?, ?)`
	res, err := s.DB.Exec(query, category.Name, category.Color)
	if err != nil {
		log.Printf("Error creating category '%s': %v", category.Name, err)
		return fmt.Errorf("failed to insert category: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID for category '%s': %v", category.Name, err)
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	category.ID = id
	log.Printf("Successfully created category '%s' with ID %d", category.Name, category.ID)
	return nil
}

func (s *sqlStore) GetCategoryByID(id int64) (*Category, error) {
	var category Category
	err := s.DB.Get(&category, "SELECT id, name, color FROM categories WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Category with ID %d not found: %v", id, err)
			return nil, nil
		}
		log.Printf("Error getting category by ID %d: %v", id, err)
		return nil, err
	}
	return &category, nil
}

func (s *sqlStore) UpdateCategory(category *Category) error {
	if category.ID == 0 {
		return fmt.Errorf("category ID cannot be zero for update")
	}
	if category.Name == "" {
		return fmt.Errorf("category name cannot be empty for update")
	}
	if category.Color == "" {
		return fmt.Errorf("category color cannot be empty for update")
	}

	query := `UPDATE categories SET name = ?, color = ? WHERE id = ?`
	res, err := s.DB.Exec(query, category.Name, category.Color, category.ID)
	if err != nil {
		log.Printf("Error updating category ID %d: %v", category.ID, err)
		return fmt.Errorf("failed to update category: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for update category ID %d: %v", category.ID, err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("No category found with ID %d to update, or data was the same.", category.ID)
		return sql.ErrNoRows
	}
	log.Printf("Successfully updated category ID %d", category.ID)
	return nil
}

func (s *sqlStore) DeleteCategory(id int64) error {
	if id == 0 {
		return fmt.Errorf("category ID cannot be zero for delete")
	}

	query := `DELETE FROM categories WHERE id = ?`
	res, err := s.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting category ID %d: %v", id, err)
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for delete category ID %d: %v", id, err)
		return fmt.Errorf("failed to get rows affected on delete: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("No category found with ID %d to delete.", id)
		return sql.ErrNoRows
	}

	log.Printf("Successfully deleted category ID %d", id)
	return nil
}
