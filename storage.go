package main

import (
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	CreateItem(x *Item) error

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

type PriceView struct {
	Sum            float64  `json:"sum,omitempty"`
	Quantity       float64  `json:"quantity,omitempty"`
	AvgPricePerOne float64  `json:"avg_price_per_one,omitempty"`
	ItemName       string   `json:"item_name,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

func (s Storage) GetPriceView(tags []string, start, end string) ([]PriceView, error) {

	resp := []PriceView{}

	stmt, args, err := sqlx.In(`SELECT 
	SUM(invoices.price_total) AS sum,
	COUNT(invoices.quantity)  AS quantity,
	AVG(invoices.price_total) AS avg_price_per_1,
	i.name                    AS item_name,
	(
		SELECT JSON_ARRAYAGG(tag)
		FROM item_tag
		WHERE tag in (?)
			AND item_tag.item_id = i.id
	)                         as tags
 FROM invoices
		  LEFT JOIN items i ON item_id = i.id
 GROUP BY i.id`, tags)
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
			addition   PriceView
			mergedTags string
		)

		err := rows.Scan(
			&addition.Sum,
			&addition.Quantity,
			&addition.AvgPricePerOne,
			&addition.ItemName,
			&mergedTags,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(mergedTags), &addition.Tags)
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
