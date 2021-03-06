package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

type Phone struct {
	ID     int
	Number string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}

	return nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	stmt := `INSERT INTO phone_numbers(value) VALUES ($1) RETURNING id`

	var id int
	err := db.QueryRow(stmt, phone).Scan(&id)

	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query(`SELECT id, value FROM phone_numbers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Phone

	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	// err := db.QueryRow(`SELECT value FROM phone_numbers WHERE id=$1`, id).Scan(&number)
	err := db.db.QueryRow(`SELECT * FROM phone_numbers WHERE value=$1`, number).Scan(&p.ID, &p.Number)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (db *DB) DeletePhone(id int) error {
	stmt := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.db.Exec(stmt, id)
	return err
}

func (db *DB) UpdatePhone(p *Phone) error {
	stmt := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.db.Exec(stmt, p.ID, p.Number)
	return err
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	if err := createPhoneNumbersTable(db); err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`
	_, err := db.Exec(statement)
	return err
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)

	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)

	if err != nil {
		return err
	}
	return nil
}

// Not using this
func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	// err := db.QueryRow(`SELECT value FROM phone_numbers WHERE id=$1`, id).Scan(&number)
	err := db.QueryRow(`SELECT * FROM phone_numbers WHERE id=$1`, id).Scan(&id, &number)

	if err != nil {
		return "", err
	}

	return number, nil
}
