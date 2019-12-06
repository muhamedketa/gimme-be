package main

import (
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	GetInvoices(start, end string) ([]Invoice, error)
	AddInvoiceSale(invoiceID int, sale Sale) error

	Close() error
}

type Storage struct {
	DB *sqlx.DB
}

func NewStorage() (StorageI, error) {
	db, err := sqlx.Open("mysql", "o3l0lijfnmr47tkw:rusk8zr7gqfhlo79@tcp(q2gen47hi68k1yrb.chr7pe7iynqr.eu-west-1.rds.amazonaws.com:3306)/c7uc48qtojxnq2gu")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return Storage{db}, err
}

// GetInvoices gets the items and merged tags from db; used semi-frequently
func (s Storage) GetInvoices(start, end string) ([]Invoice, error) {

	var resp []Invoice

	rows, err := s.DB.Query(`SELECT id                                                                              AS id,
	item_name                                                                       AS item_name,
	json_arrayagg(json_object('date', sale.date, 'total', IFNULL(sale.total, 0))) AS sales
FROM invoices i
	  LEFT JOIN invoice_sale sale on i.id = sale.invoice_id AND sale.date BETWEEN ? AND ?
GROUP BY id`, start, end)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			addition    Invoice
			mergedSales string
		)

		err := rows.Scan(
			&addition.ID,
			&addition.ItemName,
			&mergedSales,
		)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal([]byte(mergedSales), &addition.Sales); err != nil {
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

func (s Storage) AddInvoiceSale(invoiceID int, sale Sale) error {

	_, err := s.DB.Exec(`INSERT INTO invoice_sale (invoice_id, total, date) VALUES (?, ?, ?)`,
		invoiceID,
		sale.Total,
		sale.Date,
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
