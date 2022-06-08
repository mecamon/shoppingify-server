package repositories

// TODO add integration build tag

import (
	"database/sql"
	"github.com/mecamon/shoppingify-server/config"
	"github.com/mecamon/shoppingify-server/db"
	"log"
	"os"
	"testing"
)

var (
	authRepo       AuthRepo
	categoriesRepo CategoriesRepo
	itemsRepo      ItemsRepo
	listsRepo      ListsRepo
	conn           *sql.DB
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	config.Set()
	conf := config.Get()
	conn, _ = db.InitDB(conf)
	authRepo = initAuthRepo(conn, conf)
	categoriesRepo = initCategoriesRepo(conn, conf)
	itemsRepo = initItemsRepo(conn, conf)
	listsRepo = initListsRepo(conn, conf)
}

func shutdown() {
	err := conn.Close()
	if err != nil {
		log.Println("error shutting down db connection: ", err.Error())
	}
}
