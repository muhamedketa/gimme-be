package main

import "fmt"

type Invoice struct {
	ID       int     `json:"id,omitempty"`
	ItemID   int     `json:"item_id,omitempty"`
	Quantity float64 `json:"quantity,omitempty"`
	Price    float64 `json:"price,omitempty"`
}

func (x Invoice) totalPrice() float64 {
	return x.Quantity * x.Price
}

func (x Invoice) Validate() error {

	if x.ItemID == 0 || x.Quantity == 0 {
		return fmt.Errorf("missing data")
	}

	return nil

}
