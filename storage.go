package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type StorageI interface {
	CreateItem(x *Item) error

	CreateInvoice(x *Invoice) error

	Close() error
}

type Storage struct {
	DB *sql.DB
}

func NewStorage() (StorageI, error) {

	db, err := sql.Open("mysql", "root:mypassword@tcp(localhost:6603)/primary_db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return Storage{db}, err
}

func (s Storage) CreateItem(x *Item) error {

	_, err := s.DB.Exec(
		"INSERT INTO items(name) values(?)",
		x.Name,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s Storage) CreateInvoice(x *Invoice) error {
	_, err := s.DB.Exec(
		`INSERT INTO invoices(
			item_id, 
			quantity, 
			price_per_1, 
			price_total, 
			created_at
		) values(
			?,
			?,
			?,
			?,
			NOW())`,
		x.ItemID,
		x.Quantity,
		x.Price,
		x.TotalPrice(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s Storage) Close() error {
	return s.DB.Close()
}

var _ StorageI = Storage{}
