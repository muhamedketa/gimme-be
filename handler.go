package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/muhamedketa/gimme-back/api"
)

func nilFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (s server) SetRoutes() {

	s.Router.DELETE("/item/:item", nilFunc)

	s.Router.GET("/invoices", s.GetInvoicesHandler())
	s.Router.POST("/invoices", s.AddSaleHandler())

	s.Router.GET("/health", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) { w.Write([]byte("ok")) })

}

func (s server) GetInvoicesHandler() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		start, end := r.FormValue("start"), r.FormValue("end")

		invoices, err := s.Storage.GetInvoices(start, end)
		if err != nil {
			log.Println("handler", "GetInvoicesHandler", "error", err)
			api.WriteErr(w, 500, err)
			return
		}
		api.Write(w, invoices)

	}

}

func (s server) AddSaleHandler() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		fmt.Println("endpoint hit")

		data := struct {
			InvoiceID int     `json:"invoice_id,omitempty"`
			Date      string  `json:"date,omitempty"`
			Total     float64 `json:"total,omitempty"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			api.WriteErr(w, 500, err)
			return
		}
		defer r.Body.Close()

		err = s.Storage.AddInvoiceSale(data.InvoiceID, Sale{Total: data.Total, Date: data.Date})
		if err != nil {
			api.WriteErr(w, 500, err)
			return
		}

		api.Write(w, struct{}{})

	}
}
