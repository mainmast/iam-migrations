package migration

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	//Postgres migration driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	//migration file
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//Postgres general driver
	_ "github.com/lib/pq"
)

//Migrate ...
func Migrate(action string, schema string) bool {

	if strings.TrimSpace(action) == "" || strings.TrimSpace(schema) == "" {

		fmt.Println("Error with action or shcema")
		return false
	}

	fmt.Println("Selected action: " + action + "!")

	if action == "upgrade" {

		return upgrade(schema)

	} else if action == "downgrade" {

		return downgrade(schema)
	}

	return false
}

func upgrade(schema string) bool {

	db, err := sql.Open("postgres", os.Getenv("IAM_DB_URI"))

	if err != nil {

		fmt.Println("could not connect with the database")
		return false
	}

	defer db.Close()

	if _, err := db.Exec("CREATE SCHEMA " + schema); err != nil {

		fmt.Println("Error creating new schema", err)
		return false
	}

	dbURI := os.Getenv("IAM_DB_URI") + "&search_path=" + schema

	m, err := migrate.New(
		os.Getenv("CUSTOMER_MIGRATION_FILES"),
		dbURI)

	if err != nil {
		db.Exec("DROP SCHEMA " + schema + " CASCADE")
		fmt.Println("Error starting migration", err)
		return false
	}

	if err := m.Up(); err != nil {

		db.Exec("DROP SCHEMA " + schema + " CASCADE")
		fmt.Println("Error", err)
		return false
	}

	fmt.Println("upgrade run successfuly")
	return true
}

func downgrade(schema string) bool {

	dbURI := os.Getenv("IAM_DB_URI") + "&search_path=" + schema

	m, err := migrate.New(
		os.Getenv("CUSTOMER_MIGRATION_FILES"),
		dbURI)

	if err != nil {
		fmt.Println("Error starting migration", err)
		return false
	}

	if err := m.Down(); err != nil {

		fmt.Println("Error", err)
		return false

	}

	db, err := sql.Open("postgres", os.Getenv("IAM_DB_URI"))

	if err != nil {

		fmt.Println("could not connect with the database")
		return false
	}

	defer db.Close()

	if _, err := db.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {

		fmt.Println("Error deleting the schema", err)
		return false
	}

	fmt.Println("downgrading run successfuly")
	return true
}
