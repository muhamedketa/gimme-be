package main

import (
	"encoding/json"
	"gimme-back/api"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func nilFunc(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (s server) SetRoutes() {

	s.Router.POST("/item", s.CreateItem())
	s.Router.DELETE("/item/:item", nilFunc)

	s.Router.POST("/invoice", s.CreateInvoice())
	s.Router.DELETE("/invoice", nilFunc)
	s.Router.GET("/invoices", nilFunc)

	s.Router.GET("/health", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) { w.Write([]byte("ok")) })

}

func (s server) CreateItem() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var x Item
		err := json.NewDecoder(r.Body).Decode(&x)
		if err != nil {
			api.WriteErr(w, http.StatusBadRequest, err, "failed decoding")
			return
		}

		defer r.Body.Close()

		err = x.Validate()
		if err != nil {
			api.WriteErr(w, http.StatusInternalServerError, err, "failed storing")
			return
		}

		err = s.Storage.CreateItem(&x)
		if err != nil {
			log.Println(err)
			api.WriteErr(w, http.StatusInternalServerError, err, "failed storing")
			return
		}

		api.Write(w, nil)
	}
}

func (s server) CreateInvoice() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var x Invoice
		err := json.NewDecoder(r.Body).Decode(&x)
		if err != nil {
			api.WriteErr(w, http.StatusBadRequest, err, "failed decoding")
			return
		}

		defer r.Body.Close()

		err = x.Validate()
		if err != nil {
			api.WriteErr(w, http.StatusInternalServerError, err, "failed storing")
			return
		}

		err = s.Storage.CreateInvoice(&x)
		if err != nil {
			log.Println(err)
			api.WriteErr(w, http.StatusInternalServerError, err, "failed storing")
			return
		}

		api.Write(w, nil)

	}
}
