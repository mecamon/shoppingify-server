package repositories

import (
	"context"
	"database/sql"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type CategoriesRepoPostgres struct {
	Conn *sql.DB
	App  *config.App
}

func initCategoriesRepo(conn *sql.DB, app *config.App) CategoriesRepo {
	return &CategoriesRepoPostgres{
		Conn: conn,
		App:  app,
	}
}

func (r *CategoriesRepoPostgres) RegisterCategory(cat models.Category) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ID int64
	query := `INSERT INTO categories (name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING ID`
	err := r.Conn.QueryRowContext(ctx, query, cat.Name, cat.UserID, cat.CreatedAt, cat.UpdatedAt).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (r *CategoriesRepoPostgres) SearchCategoryByName(q string, skip, take int, userID int64) ([]models.CategoryDTO, error) {
	var categories []models.CategoryDTO

	query := `SELECT id, name FROM categories AS c WHERE c.user_id=$4 AND c.name ILIKE $1 OFFSET $2 LIMIT $3`
	stmt, _ := r.Conn.Prepare(query)
	defer stmt.Close()

	rows, err := stmt.Query("%"+q+"%", skip, take, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cat models.CategoryDTO
		err := rows.Scan(&cat.ID, &cat.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepoPostgres) GetAll(take, skip int, userID int64) ([]models.CategoryDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var categories []models.CategoryDTO

	query := `SELECT id, name from categories WHERE user_id = $3 ORDER BY name DESC LIMIT $1 OFFSET $2`
	rows, err := r.Conn.QueryContext(ctx, query, take, skip, userID)
	defer rows.Close()
	if err != nil {
		return categories, err
	}

	for rows.Next() {
		var category models.CategoryDTO
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}

	if rows.Err() != nil {
		return categories, rows.Err()
	}

	return categories, nil
}

func (r *CategoriesRepoPostgres) GetAllWithItemName(q string, take, skip int, userID int64) ([]models.CategoryDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var categories []models.CategoryDTO

	query := `
		SELECT c.id, c.name 
		FROM categories 
		AS c INNER JOIN items 
		ON c.id=items.category_id 
		WHERE c.user_id=$4 AND items.name 
	    ILIKE $1 LIMIT $2 OFFSET $3
	    `
	rows, err := r.Conn.QueryContext(ctx, query, "%"+q+"%", take, skip, userID)
	defer rows.Close()
	if err != nil {
		return categories, err
	}

	for rows.Next() {
		var category models.CategoryDTO
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return categories, nil
		}
		categories = append(categories, category)
	}

	if rows.Err() != nil {
		return categories, rows.Err()
	}
	return categories, nil
}

func (r *CategoriesRepoPostgres) Count(userID int64, filter ...string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int64
	var query string

	if len(filter) == 0 {
		query = `SELECT COUNT(*) FROM categories WHERE user_id=$1`
		err := r.Conn.QueryRowContext(ctx, query, userID).Scan(&count)
		if err != nil {
			return 0, err
		}
	} else {
		query = `SELECT COUNT(*) FROM categories AS c WHERE c.name LIKE $1 AND c.user_id=$2`
		err := r.Conn.QueryRowContext(ctx, query, "%"+filter[0]+"%", userID).Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}
