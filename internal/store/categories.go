package store

import (
	"database/sql" // For sql.ErrNoRows
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
)

// Category struct (as defined in models.go)

// GetAllCategories retrieves all categories from the database, ordered by name.
// It returns an empty slice if no categories are found.
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

// CreateCategory inserts a new category into the database.
// It requires Name and Color to be set on the Category struct.
// The ID of the newly created category is set on the input Category struct.
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

// GetCategoryByID retrieves a single category by its ID.
// It returns nil if the category is not found.
func (s *sqlStore) GetCategoryByID(id int64) (*Category, error) {
	var category Category
	err := s.DB.Get(&category, "SELECT id, name, color FROM categories WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Category with ID %d not found: %v", id, err)
			return nil, nil // Or a specific "not found" error
		}
		log.Printf("Error getting category by ID %d: %v", id, err)
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category in the database.
// It requires ID, Name, and Color to be set on the Category struct.
// Returns sql.ErrNoRows if no category with the given ID is found or if data was the same.
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
		return sql.ErrNoRows // Use sql.ErrNoRows to indicate not found or no change
	}
	log.Printf("Successfully updated category ID %d", category.ID)
	return nil
}

// DeleteCategory removes a category from the database by its ID.
// Returns sql.ErrNoRows if no category with the given ID is found.
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
		return sql.ErrNoRows // Use sql.ErrNoRows to indicate not found
	}

	log.Printf("Successfully deleted category ID %d", id)
	return nil
}
