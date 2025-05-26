package store

import (
	"database/sql" // For sql.ErrNoRows
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
)

// Category struct (as defined in models.go)

// GetAllCategories ... (as before)
func GetAllCategories(db *sqlx.DB) ([]Category, error) {
	var categories []Category
	err := db.Select(&categories, "SELECT id, name, color FROM categories ORDER BY name ASC")
	if err != nil {
		log.Printf("Error getting all categories: %v", err)
		return nil, err
	}
	if categories == nil {
		categories = []Category{} 
	}
	return categories, nil
}

// CreateCategory ... (as before)
func CreateCategory(db *sqlx.DB, category *Category) error {
	if category.Name == "" { return fmt.Errorf("category name cannot be empty") }
	if category.Color == "" { return fmt.Errorf("category color cannot be empty") }
	query := `INSERT INTO categories (name, color) VALUES (?, ?)`
	res, err := db.Exec(query, category.Name, category.Color)
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
func GetCategoryByID(db *sqlx.DB, id int64) (*Category, error) {
	var category Category
	err := db.Get(&category, "SELECT id, name, color FROM categories WHERE id = ?", id)
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
// It ensures the ID in the category struct is used for the WHERE clause.
func UpdateCategory(db *sqlx.DB, category *Category) error {
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
	res, err := db.Exec(query, category.Name, category.Color, category.ID)
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
func DeleteCategory(db *sqlx.DB, id int64) error {
	if id == 0 {
		return fmt.Errorf("category ID cannot be zero for delete")
	}

	query := `DELETE FROM categories WHERE id = ?`
	res, err := db.Exec(query, id)
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
