package main

import (
	"database/sql"
	"log"

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

	ch := http.NewCustomerHandler(customerManager)
	ca := http.NewCheckingAccountHandler(checkinAccountRepoManager)

	router := mux.NewRouter()
	userRouteCreate := router.NewRoute()
	userRouteCreate.Path("/customer").HandlerFunc(ch.CreateCustomer).Methods("POST")
	userRouteUpdate := router.NewRoute()
	userRouteUpdate.Path("/customer/{id}/").HandlerFunc(ch.UpdateCustomer).Methods("PUT")
	userRouteGet := router.NewRoute()
	userRouteGet.Path("/customer/{id}/").HandlerFunc(ch.GetCustomer).Methods("GET")

	checkAccRouteCreate := router.NewRoute()
	checkAccRouteCreate.Path("/checking-account").HandlerFunc(ca.CreateCheckingAccount).Methods("POST")
	checkAccRouteUpdate := router.NewRoute()
	checkAccRouteUpdate.Path("/checking-account/{id}/").HandlerFunc(ca.UpdateCheckingAccount).Methods("PUT")
	checkAccRouteGet := router.NewRoute()
	checkAccRouteGet.Path("/checking-account/{id}/").HandlerFunc(ca.GetCheckingAccount).Methods("GET")
	panic(gohttp.ListenAndServe(":8080", router))
}
