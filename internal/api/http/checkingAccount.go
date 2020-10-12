package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/gorilla/mux"

	"github.com/fwiedmann/ezb/domain/entity/checking_account"
)

type CheckingAccount struct {
	Number         string `json:"number"`
	Name           string `json:"name"`
	OverDraftLimit string `json:"over_draft_limit"`
	Balance        string `json:"balance"`
	Pin            string `json:"pin"`
}

func NewCheckingAccountHandler(m checking_account.Manager) *CheckingAccountHandler {
	return &CheckingAccountHandler{manager: m}
}

type CheckingAccountHandler struct {
	manager checking_account.Manager
}

func (c *CheckingAccountHandler) CreateCheckingAccount(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	var requestCheckingAccount CheckingAccount
	if err := json.Unmarshal(body, &requestCheckingAccount); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	if requestCheckingAccount.OverDraftLimit == "" {
		requestCheckingAccount.OverDraftLimit = "0"
	}
	parsedOverdraftLimit, err := strconv.ParseFloat(requestCheckingAccount.OverDraftLimit, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	id, err := c.manager.Create(r.Context(), checking_account.CheckingAccount{
		Name:           requestCheckingAccount.Name,
		OverDraftLimit: parsedOverdraftLimit,
	}, requestCheckingAccount.Pin)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	requestCheckingAccount.Number = id.String()
	resp, err := json.Marshal(&requestCheckingAccount)
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

func (c *CheckingAccountHandler) UpdateCheckingAccount(w http.ResponseWriter, r *http.Request) {
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

	var requestCheckingAccount CheckingAccount
	if err := json.Unmarshal(body, &requestCheckingAccount); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	parsedOverdraftLimit, err := strconv.ParseFloat(requestCheckingAccount.OverDraftLimit, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	parsedNumber, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	err = c.manager.Update(r.Context(), checking_account.CheckingAccount{
		Number:         parsedNumber,
		Name:           requestCheckingAccount.Name,
		OverDraftLimit: parsedOverdraftLimit,
	}, requestCheckingAccount.Pin)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	requestCheckingAccount.Number = parsedNumber.String()
	resp, err := json.Marshal(&requestCheckingAccount)
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

func (c *CheckingAccountHandler) GetCheckingAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	parsedNumber, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}
	ca, err := c.manager.Get(r.Context(), parsedNumber)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), 422)
		return
	}

	respCheckingAccount := CheckingAccount{
		Number:         ca.Number.String(),
		Name:           ca.Name,
		OverDraftLimit: fmt.Sprintf("%b", ca.OverDraftLimit),
		Balance:        fmt.Sprintf("%b", ca.GetCurrentBalance()),
	}

	resp, err := json.Marshal(&respCheckingAccount)
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
