package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type ItemsRepoPostgres struct {
	Conn *sql.DB
	App  *config.App
}

func initItemsRepo(conn *sql.DB, app *config.App) ItemsRepo {
	return &ItemsRepoPostgres{
		Conn: conn,
		App:  app,
	}
}

func (i *ItemsRepoPostgres) Register(item models.Item) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	query := `INSERT INTO items (
		   name, 
		   note, 
		   category_id, 
		   image_url, 
		   is_active,
		   created_at, 
		   updated_at
	    ) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING ID`

	var ID int64

	err := i.Conn.QueryRowContext(ctx, query,
		item.Name,
		item.Note,
		item.CategoryID,
		item.ImageURL,
		item.IsActive,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (i *ItemsRepoPostgres) GetAll(take, skip int) ([]models.ItemDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	stmt := `
		SELECT i.id, i.name, c.id 
		FROM items as i
		INNER JOIN categories c
		ON i.category_id=c.id
		WHERE i.is_active=true 
		LIMIT $1 OFFSET $2`
	rows, err := i.Conn.QueryContext(ctx, stmt, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.ItemDTO

	for rows.Next() {
		var item models.ItemDTO
		err := rows.Scan(&item.ID, &item.Name, &item.CategoryID)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (i *ItemsRepoPostgres) GetByID(itemID int64) (models.ItemDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var item models.ItemDTO

	stmt := `
		SELECT i.id, i.name, c.id 
		FROM items AS i 
		INNER JOIN categories c
		ON i.category_id=c.id
		WHERE i.id=$1`
	err := i.Conn.QueryRowContext(ctx, stmt, itemID).Scan(&item.ID, &item.Name, &item.CategoryID)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (i *ItemsRepoPostgres) GetAllByCategoryID(categoryID int64) ([]models.ItemDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var items []models.ItemDTO

	query := `
		SELECT i.id, i.name, c.id 
		FROM items AS i
		INNER JOIN categories c
		ON i.category_id=c.id
		WHERE i.category_id=$1 AND is_active=true`
	rows, err := i.Conn.QueryContext(ctx, query, categoryID)
	defer rows.Close()
	if err != nil {
		return items, err
	}

	for rows.Next() {
		var item models.ItemDTO
		err := rows.Scan(&item.ID, &item.Name, &item.CategoryID)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return items, rows.Err()
	}
	return items, nil
}

func (i *ItemsRepoPostgres) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int64

	query := `SELECT COUNT(*) FROM items WHERE is_active=true`
	err := i.Conn.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (i *ItemsRepoPostgres) GetDetails(id int64) (models.ItemDetailedDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var item models.ItemDetailedDTO

	query := `
		SELECT i.id, i.name, c.name, i.note, i.image_url 
		FROM items AS i
		INNER JOIN categories c
		ON i.category_id=c.id
		WHERE i.id=$1`
	err := i.Conn.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.CategoryName,
		&item.Note,
		&item.ImageURL)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (i *ItemsRepoPostgres) Disable(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE items SET is_active=false WHERE id=$1`
	result, err := i.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("0 rows affected")
	}
	return nil
}
