package repositories

import (
	"database/sql"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/models"
	"github.com/mecamon/shoppingify-server/utils"
)

type AuthRepoPostgres struct {
	Conn *sql.DB
	App  *config.App
}

func initAuthRepo(conn *sql.DB, app *config.App) AuthRepo {
	return &AuthRepoPostgres{Conn: conn, App: app}
}

func (r *AuthRepoPostgres) Register(user models.User) error {
	query := `INSERT INTO users (
       name, 
	   lastname, 
	   email, 
	   password, 
	   is_active, 
	   is_visitor, 
	   login_code, 
	   created_at, 
	   updated_at
	   ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.Conn.Exec(
		query, user.Name,
		user.Lastname,
		user.Email,
		user.Password,
		user.IsActive,
		user.IsVisitor,
		user.LoginCode,
		user.CreatedAt,
		user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepoPostgres) SearchUserByEmail(email string) (models.User, error) {
	user := models.User{}
	query := `SELECT name, lastname, email FROM users WHERE email=$1`
	row := r.Conn.QueryRow(query, email)

	err := row.Scan(&user.Name, &user.Lastname, &user.Email)
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *AuthRepoPostgres) CheckUserPassword(email, password string) (bool, error) {
	user := models.User{}
	query := `SELECT name, lastname, email, password FROM users WHERE email=$1`
	row := r.Conn.QueryRow(query, email)

	err := row.Scan(&user.Name, &user.Lastname, &user.Email, &user.Password)
	if err != nil {
		return false, err
	}

	hasCorrectPass, _ := utils.CompareHashAndPass(user.Password, password)
	return hasCorrectPass, nil
}
