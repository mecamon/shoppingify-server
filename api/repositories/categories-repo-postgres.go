package repositories

import (
	"database/sql"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"log"
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

func (r CategoriesRepoPostgres) RegisterCategory(cat models.Category) (int64, error) {
	var ID int64
	query := `INSERT INTO categories (name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING ID`
	err := r.Conn.QueryRow(query, cat.Name, cat.UserID, cat.CreatedAt, cat.UpdatedAt).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (r CategoriesRepoPostgres) SearchCategoryByName(q string, skip, take int) ([]models.CategoryDTO, error) {
	var categories []models.CategoryDTO

	query := `SELECT id, name FROM categories AS c WHERE c.name LIKE $1 OFFSET $2 LIMIT $3`
	stmt, _ := r.Conn.Prepare(query)
	defer stmt.Close()

	rows, err := stmt.Query("%"+q+"%", skip, take)
	defer rows.Close()
	if err != nil {
		log.Println("ERROR:", err.Error())
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
