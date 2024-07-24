package main

import (
	"github.com/google/uuid"
	"net/http"
)

// Delete -> DELETE /company/{id} -> delete, companys.html

// Edit   -> GET /company/edit/{id} -> row-edit.html
// Save   ->   PUT /company/{id} -> update, row.html
// Cancel ->	 GET /company/{id} -> nothing, row.html

// Add    -> GET /company/add/ -> companys-add.html (target body with row-add.html and row.html)
// Save   ->   POST /company -> add, companys.html (target body without row-add.html)
// Cancel ->	 GET /company -> nothing, companys.html

func index(r *http.Request) *Response {
	return HTML(http.StatusOK, html, "index.html", data, nil)
}

// GET /company/add
func companyAdd(r *http.Request) *Response {
	return HTML(http.StatusOK, html, "company-add.html", data, nil)
}

// /GET company/edit/{id}
func companyEdit(r *http.Request) *Response {
	id, _ := PathLast(r)
	uuid := uuid.MustParse(id)
	row := getCompanyByID(uuid)
	return HTML(http.StatusOK, html, "row-edit.html", row, nil)
}

// GET /company
// GET /company/{id}
// DELETE /company/{id}
// PUT /company/{id}
// POST /company
func companies(r *http.Request) *Response {
	switch r.Method {

	case http.MethodDelete:
		id, _ := PathLast(r)
		uuid := uuid.MustParse(id)
		deleteCompany(uuid)
		return HTML(http.StatusOK, html, "companies.html", data, nil)

	//cancel
	case http.MethodGet:
		id, segments := PathLast(r)
		if segments > 1 {
			//cancel edit
			uuid := uuid.MustParse(id)
			row := getCompanyByID(uuid)
			return HTML(http.StatusOK, html, "row.html", row, nil)
		} else {
			//cancel add
			return HTML(http.StatusOK, html, "companies.html", data, nil)
		}

	//save edit
	case http.MethodPut:
		id, _ := PathLast(r)
		uuid := uuid.MustParse(id)
		row := getCompanyByID(uuid)
		r.ParseForm()
		row.Company = r.Form.Get("company")
		row.Contact = r.Form.Get("contact")
		row.Country = r.Form.Get("country")
		updateCompany(row)
		return HTML(http.StatusOK, html, "row.html", row, nil)

	//save add
	case http.MethodPost:
		row := Company{}
		r.ParseForm()
		row.Company = r.Form.Get("company")
		row.Contact = r.Form.Get("contact")
		row.Country = r.Form.Get("country")
		addCompany(row)
		return HTML(http.StatusOK, html, "companies.html", data, nil)
	}

	return Empty(http.StatusNotImplemented)
}
