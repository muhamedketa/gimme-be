package main

type Invoice struct {
	ID       int    `json:"id,omitempty"`
	ItemName string `json:"item_name,omitempty"`
	Sales    []Sale `json:"sales,omitempty"`
}

// Sale gives the data on total sale on an item
type Sale struct {
	Date     string  `json:"date,omitempty"`
	Total    float64 `json:"total,omitempty"`
	Quantity float64 `json:"quantity,omitempty"`
}

func (x Sale) priceOne() float64 {
	return x.Total / x.Quantity
}
