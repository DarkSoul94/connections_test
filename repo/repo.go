package repo

import (
	"database/sql"
	"fmt"

	"github.com/DarkSoul94/connections_test/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	"github.com/oklog/ulid"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbName   = "postgres"
)

func ConnectDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return conn
}

func RunMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName,
		driver)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		fmt.Println(err)
	}
}

func CreateUser(db *sql.DB, user models.User) error {
	query := `INSERT INTO users(id, name, age) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, user.ID.String(), user.Name, user.Age)

	return err
}

func GetUser(db *sql.DB, id ulid.ULID) (models.User, error) {
	var user models.User
	var idStr string
	query := `SELECT id, name, age FROM users WHERE id = $1`

	row := db.QueryRow(query, id.String())
	err := row.Scan(&idStr, &user.Name, &user.Age)
	if err != nil {
		return models.User{}, err
	}

	user.ID = ulid.MustParse(idStr)

	return user, nil
}

func DeleteUser(db *sql.DB, id ulid.ULID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(query, id.String())

	return err
}
