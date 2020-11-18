package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/fwiedmann/ezb/domain/usecase/assign_customer_to_checking_account"
)

type Mapping struct {
	CustomerID            string `json:"customer_id"`
	CheckingAccountNumber string `json:"checking_account_number"`
}

type customerCheckingAccountMapper struct {
	mapper *assign_customer_to_checking_account.UseCase
}

func (cma *customerCheckingAccountMapper) CreateMapping(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	var mappingRequest Mapping
	if err := json.Unmarshal(body, &mappingRequest); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	parsedCustomerID, err := uuid.Parse(mappingRequest.CustomerID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	parsedCheckingAccountNumber, err := uuid.Parse(mappingRequest.CheckingAccountNumber)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	if err := cma.mapper.CreateMapping(r.Context(), parsedCustomerID, parsedCheckingAccountNumber); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
