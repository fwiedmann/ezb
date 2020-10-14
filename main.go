package main

import (
	"database/sql"
	"log"

	"github.com/fwiedmann/ezb/domain/entity/customer_checking_account_mapping"
	"github.com/fwiedmann/ezb/domain/usecase/assign_customer_to_checking_account"

	"github.com/fwiedmann/ezb/domain/usecase/debit"

	"github.com/fwiedmann/ezb/domain/usecase/deposit"

	"github.com/fwiedmann/ezb/domain/entity/checking_account"

	"github.com/fwiedmann/ezb/internal/api/http"

	"github.com/gorilla/mux"

	gohttp "net/http"

	"github.com/fwiedmann/ezb/domain/entity/customer"
)
import _ "github.com/go-sql-driver/mysql"

func main() {
	db, err := sql.Open("mysql", "root:secureByDefault@(127.0.0.1:3306)/ezb")
	if err != nil {
		log.Fatal(err)
	}

	customerRepo, err := customer.NewMySqlRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	customerManager := customer.NewManager(customerRepo)

	checkinAccountRepo, err := checking_account.NewMySqlRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	checkinAccountRepoManager := checking_account.NewManager(checkinAccountRepo)

	depositUseCase := deposit.NewUseCase(checkinAccountRepoManager)
	debitUseCase := debit.NewUseCase(checkinAccountRepoManager)

	ch := http.NewCustomerHandler(customerManager)
	ca := http.NewCheckingAccountHandler(checkinAccountRepoManager, depositUseCase, debitUseCase)

	mapperRepo, err := customer_checking_account_mapping.NewMySqlRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	mapperHandler := http.NewCustomerCheckingAccountMapper(assign_customer_to_checking_account.NewUseCase(customer_checking_account_mapping.NewManager(mapperRepo)))

	router := mux.NewRouter()
	router.NewRoute().Path("/customer").HandlerFunc(ch.CreateCustomer).Methods("POST")
	router.NewRoute().Path("/customer/{id}/").HandlerFunc(ch.UpdateCustomer).Methods("PUT")
	router.NewRoute().Path("/customer/{id}/").HandlerFunc(ch.GetCustomer).Methods("GET")

	router.NewRoute().Path("/checking-account").HandlerFunc(ca.CreateCheckingAccount).Methods("POST")
	router.NewRoute().Path("/checking-account/{id}/").HandlerFunc(ca.UpdateCheckingAccount).Methods("PUT")
	router.NewRoute().Path("/checking-account/{id}/").HandlerFunc(ca.GetCheckingAccount).Methods("GET")
	router.NewRoute().Path("/checking-account/{id}/deposit").HandlerFunc(ca.Deposit).Methods("POST")
	router.NewRoute().Path("/checking-account/{id}/debit").HandlerFunc(ca.Debit).Methods("POST")

	mapperRoutePost := router.NewRoute()
	mapperRoutePost.Path("/customer-checking-account-mapping").HandlerFunc(mapperHandler.CreateMapping).Methods("POST")
	panic(gohttp.ListenAndServe(":8080", router))
}
