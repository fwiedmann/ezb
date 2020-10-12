package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"

	"github.com/fwiedmann/ezb/domain/entity/customer"
)

type Customer struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"`
}

func NewCustomerHandler(m customer.Manager) *CustomerHandler {
	return &CustomerHandler{manager: m}
}

type CustomerHandler struct {
	manager customer.Manager
}

func (c *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	var requestCustomer Customer
	if err := json.Unmarshal(body, &requestCustomer); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	id, err := c.manager.Create(r.Context(), customer.Customer{
		FirstName: requestCustomer.FirstName,
		LastName:  requestCustomer.LastName,
		Gender:    requestCustomer.Gender,
		Birthdate: requestCustomer.Birthdate,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	requestCustomer.ID = id.String()
	resp, err := json.Marshal(&requestCustomer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	_, err = fmt.Fprint(w, string(resp))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	var requestCustomer Customer
	if err := json.Unmarshal(body, &requestCustomer); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	err = c.manager.Update(r.Context(), customer.Customer{
		ID:        parsedID,
		FirstName: requestCustomer.FirstName,
		LastName:  requestCustomer.LastName,
		Gender:    requestCustomer.Gender,
		Birthdate: requestCustomer.Birthdate,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	requestCustomer.ID = parsedID.String()
	resp, err := json.Marshal(&requestCustomer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	_, err = fmt.Fprint(w, string(resp))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	cu, err := c.manager.Get(r.Context(), parsedID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), 404)
		return
	}

	respCustomer := Customer{
		ID:        cu.ID.String(),
		FirstName: cu.FirstName,
		LastName:  cu.LastName,
		Gender:    cu.Gender,
		Birthdate: cu.Birthdate,
	}

	resp, err := json.Marshal(&respCustomer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	_, err = fmt.Fprint(w, string(resp))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	w.WriteHeader(http.StatusOK)
}
