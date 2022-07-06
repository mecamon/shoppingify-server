package repositories

import (
	"context"
	"database/sql"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"time"
)

type ItemsSummaryRepoImpl struct {
	Conn *sql.DB
	App  *config.App
}

func initItemsSummaryRepoImpl(conn *sql.DB, app *config.App) ItemsSummaryRepo {
	return &ItemsSummaryRepoImpl{
		Conn: conn,
		App:  app,
	}
}

func (i ItemsSummaryRepoImpl) Add(userID int64, itemsInfo models.ItemsSummary) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `
		INSERT INTO items_summary(user_id, month, year, quantity)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (user_id, year, month) DO
		UPDATE SET quantity = $4 WHERE items_summary.user_id = $1 AND items_summary.month = $2 AND items_summary.year = $3
	`

	_, err := i.Conn.ExecContext(ctx, query, userID, itemsInfo.Month, itemsInfo.Year, itemsInfo.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (i ItemsSummaryRepoImpl) GetMonthly(userID int64, year int64) (models.ItemsSummaryByMonthDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var itemsByMonth models.ItemsSummaryByMonthDTO
	itemsByMonth.Year = year

	query := `
		SELECT sum(i.quantity) AS quantity, i.month 
		FROM items_summary AS i 
		WHERE i.user_id=$1 AND i.year = $2 GROUP BY i.month
		`

	rows, err := i.Conn.QueryContext(ctx, query, userID, year)
	defer rows.Close()
	if err != nil {
		return itemsByMonth, err
	}
	for rows.Next() {
		var item models.ItemsSummaryByMonth
		err := rows.Scan(&item.Quantity, &item.Month)
		if err != nil {
			return itemsByMonth, err
		}
		itemsByMonth.Months = append(itemsByMonth.Months, item)
	}
	if rows.Err() != nil {
		return itemsByMonth, rows.Err()
	}
	return itemsByMonth, nil
}

func (i ItemsSummaryRepoImpl) GetYearly(userID int64) ([]models.ItemsSummaryByYearDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var itemsByYear []models.ItemsSummaryByYearDTO

	query := `SELECT sum(quantity), year FROM items_summary WHERE user_id = $1 GROUP BY year`

	rows, err := i.Conn.QueryContext(ctx, query, userID)
	defer rows.Close()
	if err != nil {
		return itemsByYear, err
	}
	for rows.Next() {
		var item models.ItemsSummaryByYearDTO
		err := rows.Scan(&item.Quantity, &item.Year)
		if err != nil {
			return itemsByYear, err
		}
		itemsByYear = append(itemsByYear, item)
	}
	if rows.Err() != nil {
		return itemsByYear, rows.Err()
	}
	return itemsByYear, nil
}
