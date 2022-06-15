package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type TopItemsRepoImpl struct {
	Conn *sql.DB
	App  *config.App
}

func initTopItemsRepo(conn *sql.DB, app *config.App) TopItemsRepo {
	return &TopItemsRepoImpl{
		Conn: conn,
		App:  app,
	}
}

func (t TopItemsRepoImpl) Add(userID, itemID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ID int64

	query := `INSERT INTO top_items (user_id, item_id, sum_quantity) VALUES ($1, $2, $3) RETURNING ID`
	err := t.Conn.QueryRowContext(ctx, query, userID, itemID, 1).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (t TopItemsRepoImpl) Update(userID, itemID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE top_items as t SET sum_quantity=sum_quantity+1 WHERE t.user_id=$1 AND t.item_id=$2`
	result, err := t.Conn.ExecContext(ctx, stmt, userID, itemID)
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

func (t TopItemsRepoImpl) GetTop(userID int64, take int) ([]models.TopItemDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var topItems []models.TopItemDTO

	query := `
		SELECT t.id, t.item_id, i.name, t.sum_quantity FROM top_items as t 
		    INNER JOIN items i ON i.id=t.item_id WHERE t.user_id=$1
			ORDER BY t.sum_quantity DESC 
			LIMIT $2
		`
	rows, err := t.Conn.QueryContext(ctx, query, userID, take)
	if err != nil {
		return topItems, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.TopItemDTO
		err := rows.Scan(&item.ID, &item.ItemID, &item.Name, &item.SumQuantity)
		if err != nil {
			return topItems, err
		}
		topItems = append(topItems, item)
	}

	if rows.Err() != nil {
		return topItems, rows.Err()
	}
	return topItems, nil
}

func (t *TopItemsRepoImpl) GetAll(userID int64) ([]models.TopItemDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var topItems []models.TopItemDTO

	query := `SELECT id, item_id, sum_quantity FROM top_items WHERE user_id=$1`
	rows, err := t.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		return topItems, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.TopItemDTO
		err := rows.Scan(&item.ID, &item.ItemID, &item.SumQuantity)
		if err != nil {
			return topItems, err
		}
		topItems = append(topItems, item)
	}
	if rows.Err() != nil {
		return topItems, rows.Err()
	}
	return topItems, nil
}
