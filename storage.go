package main

import (
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	CreateItem(x *Item) error

	GetAllItems() ([]Item, error)

	CreateInvoice(x *Invoice) error

	GetPriceView(tags []string, start, end string) ([]PriceView, error)

	Close() error
}

type Storage struct {
	DB *sqlx.DB
}

func NewStorage() (StorageI, error) {

	db, err := sqlx.Open("mysql", "root:mypassword@tcp(localhost:6603)/primary_db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return Storage{db}, err
}

// GetAllItems gets the items and merged tags from db; used semi-frequently
func (s Storage) GetAllItems() ([]Item, error) {

	var resp []Item

	rows, err := s.DB.Query(`SELECT items.id as id, items.name as name, json_arrayagg(it.tag) as tags
	from items
			 LEFT JOIN item_tag it on items.id = it.item_id
	GROUP BY items.id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			addition   Item
			mergedTags string
		)

		err := rows.Scan(
			&addition.ID,
			&addition.Name,
			&mergedTags,
		)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal([]byte(mergedTags), &addition.Tags); err != nil {
			return nil, err
		}

		resp = append(resp, addition)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return resp, nil

}

// PriceView is the type of data on graph to show how much was spent in a period of time
type PriceView struct {
	Sum            float64 `json:"sum,omitempty"`
	Quantity       float64 `json:"quantity,omitempty"`
	AvgPricePerOne float64 `json:"avg_price_per_one,omitempty"`
	ItemID         int     `json:"item_id,omitempty"`
}

//TODO:  use itemID on graphs, make a separate get items which also returns tags per items
// and only return itemID on views

func (s Storage) GetPriceView(tags []string, start, end string) ([]PriceView, error) {

	resp := []PriceView{}

	stmt, args, err := sqlx.In(`SELECT SUM(price_total) AS sum,
	COUNT(quantity)  AS quantity,
	AVG(price_per_1) AS avg_price_per_1,
	invoices.item_id
FROM invoices
	  INNER JOIN items i on invoices.item_id = i.id
	  INNER JOIN item_tag it on i.id = it.item_id AND it.tag in (?)
WHERE created_at BETWEEN ? AND ?
GROUP BY item_id`, tags, start, end)
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			addition PriceView
		)

		err := rows.Scan(
			&addition.Sum,
			&addition.Quantity,
			&addition.AvgPricePerOne,
			&addition.ItemID,
		)
		if err != nil {
			return nil, err
		}

		resp = append(resp, addition)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return resp, nil
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
		`INSERT INTO invoices (
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
			NOW() )`,
		x.ItemID,
		x.Quantity,
		x.Price,
		x.totalPrice(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Close implements ioCloser
func (s Storage) Close() error {
	return s.DB.Close()
}

var _ StorageI = Storage{}
