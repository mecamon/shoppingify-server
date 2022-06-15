package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"strings"
	"time"
)

type ListsRepoPostgres struct {
	Conn *sql.DB
	App  *config.App
}

func initListsRepo(conn *sql.DB, app *config.App) ListsRepo {
	return &ListsRepoPostgres{
		Conn: conn,
		App:  app,
	}
}

func (r *ListsRepoPostgres) Create(list models.List) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO lists (
			name, 
			is_completed, 
			is_cancelled, 
			user_id, 
			created_at, 
			updated_at, 
			completed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING ID`

	var insertedID int64

	err := r.Conn.QueryRowContext(ctx, query,
		list.Name,
		list.IsCompleted,
		list.IsCancelled,
		list.UserID,
		list.CreatedAt,
		list.UpdatedAt,
		list.CompletedAt,
	).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ListsRepoPostgres) UpdateActiveListName(userID int64, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE lists SET name=$1, updated_at=$2 WHERE lists.user_id=$3 AND lists.is_completed=false AND lists.is_cancelled=false`
	result, err := r.Conn.ExecContext(ctx, stmt, name, time.Now().Unix(), userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("there is no active list")
	}
	return nil
}

func (r *ListsRepoPostgres) GetActive(userID int64) (models.ListDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var listID int64
	var listName string
	var listCreatedAt int64
	list := models.ListDTO{}

	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return list, err
	}
	defer tx.Rollback()

	query1 := `SELECT l.id, l.name, l.created_at FROM lists AS l WHERE l.user_id=$1 AND l.is_completed=false AND l.is_cancelled=false`
	row := tx.QueryRowContext(ctx, query1, userID)
	err = row.Scan(&listID, &listName, &listCreatedAt)
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return list, err
	}

	var items []models.SelectedItemDTO

	query2 := `
			SELECT i_sel.id, i_sel.quantity, i.name, i.id 
			FROM items_selected AS i_sel 
			INNER JOIN items AS i ON i_sel.item_id=i.id
			WHERE i_sel.list_id=$1
    	`
	rows, err := tx.QueryContext(ctx, query2, listID)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		item := models.SelectedItemDTO{}
		err := rows.Scan(&item.ID, &item.Quantity, &item.Name, &item.ItemID)
		if err != nil {
			return list, nil
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return list, rows.Err()
	}

	if err = tx.Commit(); err != nil {
		return list, err
	}

	list.ID = listID
	list.Name = listName
	list.Date = time.Unix(listCreatedAt, 0)
	list.Items = items

	return list, nil
}

func (r *ListsRepoPostgres) AddItemToList(item models.SelectedItem) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var isCompleted, isCancelled bool

	query1 := `SELECT is_completed, is_cancelled FROM lists WHERE lists.id=$1`
	row1 := tx.QueryRowContext(ctx, query1, item.ListID)
	if err = row1.Scan(&isCompleted, &isCancelled); err != nil {
		return 0, err
	}
	if isCompleted || isCancelled {
		return 0, errors.New("cannot add item to inactive list")
	}

	var insertedID int64
	query2 := `INSERT INTO items_selected (
			item_id, 
			quantity, 
			is_completed, 
			list_id, 
			created_at, 
			updated_at
			) VALUES ($1, $2, $3, $4, $5, $6) 
		  	RETURNING ID
		  	`
	row2 := r.Conn.QueryRowContext(ctx, query2,
		item.ItemID,
		item.Quantity,
		item.IsCompleted,
		item.ListID,
		item.CreatedAt,
		item.UpdatedAt)

	if err = row2.Scan(&insertedID); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ListsRepoPostgres) IsListActive(listID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var isCompleted, isCancelled bool

	query := `SELECT is_completed, is_cancelled FROM lists WHERE lists.id=$1`
	row := r.Conn.QueryRowContext(ctx, query, listID)
	if err := row.Scan(&isCompleted, &isCancelled); err != nil {
		return false, err
	}
	if isCompleted || isCancelled {
		return false, nil
	}
	return !isCancelled && !isCompleted, nil
}

func (r *ListsRepoPostgres) DeleteItemFromList(itemSelID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM items_selected AS i_sel WHERE i_sel.id=$1`
	result, err := r.Conn.ExecContext(ctx, stmt, itemSelID)
	if err != nil {
		return err
	}
	numberOfRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numberOfRows == 0 {
		return errors.New("item does not exist")
	}
	return nil
}

func (r *ListsRepoPostgres) CompleteItemSelected(itemSelID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE items_selected SET is_completed=true WHERE items_selected.id=$1`
	result, err := r.Conn.ExecContext(ctx, stmt, itemSelID)
	if err != nil {
		return err
	}
	numberOfRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numberOfRows == 0 {
		return errors.New("item does not exist")
	}
	return nil
}

func (r *ListsRepoPostgres) CancelActive(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE lists SET is_cancelled=true WHERE user_id=$1 AND is_completed=false AND is_cancelled=false`
	result, err := r.Conn.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}
	numberOfRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numberOfRows == 0 {
		return errors.New("there is no active list")
	}
	return nil
}

func (r *ListsRepoPostgres) CompleteActive(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE lists SET is_completed=true WHERE user_id=$1 AND is_completed=false AND is_cancelled=false`
	result, err := r.Conn.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}
	numberOfRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numberOfRows == 0 {
		return errors.New("there is no active list")
	}
	return nil
}

func (r *ListsRepoPostgres) UpdateItemsSelected(items []models.UpdateSelItemDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, i := range items {
		stmt := `UPDATE items_selected AS i_sel SET quantity=$1 WHERE i_sel.id=$2`
		result, err := tx.ExecContext(ctx, stmt, i.Quantity, i.ItemID)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New(fmt.Sprintf("%d 404", i.ItemID))
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
