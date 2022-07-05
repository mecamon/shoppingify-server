package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type TopCategoriesRepoImpl struct {
	Conn *sql.DB
	App  *config.App
}

func initTopCategoriesRepo(conn *sql.DB, app *config.App) TopCategoriesRepo {
	return &TopCategoriesRepoImpl{
		Conn: conn,
		App:  app,
	}
}

func (t TopCategoriesRepoImpl) Add(userID, categoryID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ID int64

	query := `INSERT INTO top_categories (user_id, category_id, sum_quantity) VALUES ($1, $2, $3) RETURNING ID`
	err := t.Conn.QueryRowContext(ctx, query, userID, categoryID, 1).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (t TopCategoriesRepoImpl) Update(userID, categoryID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE top_categories SET sum_quantity=sum_quantity+1 WHERE user_id=$1 AND category_id=$2`
	result, err := t.Conn.ExecContext(ctx, stmt, userID, categoryID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows were affected")
	}
	return nil
}

func (t TopCategoriesRepoImpl) GetTop(userID int64, take int) ([]models.TopCategoryDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var topCategories []models.TopCategoryDTO

	query := `
		SELECT t.id, c.id, t.sum_quantity, c.name
		FROM top_categories as t 
		    INNER JOIN categories c ON c.id=t.category_id WHERE t.user_id=$1
			ORDER BY t.sum_quantity DESC 
			LIMIT $2
		`
	rows, err := t.Conn.QueryContext(ctx, query, userID, take)
	if err != nil {
		return topCategories, err
	}
	defer rows.Close()

	for rows.Next() {
		var topCat models.TopCategoryDTO
		err := rows.Scan(&topCat.ID, &topCat.CategoryID, &topCat.SumQuantity, &topCat.Name)
		if err != nil {
			return topCategories, err
		}
		topCategories = append(topCategories, topCat)
	}

	if rows.Err() != nil {
		return topCategories, rows.Err()
	}
	return topCategories, nil
}

func (t *TopCategoriesRepoImpl) GetAll(userID int64) ([]models.TopCategoryDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var topCategories []models.TopCategoryDTO

	query := `
		SELECT t.id, t.category_id, t.sum_quantity, c.name 
		FROM top_categories AS t
		INNER JOIN categories c ON t.category_id=c.id
		WHERE t.user_id=$1`
	rows, err := t.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		return topCategories, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat models.TopCategoryDTO
		err := rows.Scan(&cat.ID, &cat.CategoryID, &cat.SumQuantity, &cat.Name)
		if err != nil {
			return topCategories, err
		}
		topCategories = append(topCategories, cat)
	}
	if rows.Err() != nil {
		return topCategories, rows.Err()
	}
	return topCategories, nil
}
