package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fwiedmann/ezb/domain/usecase/assign_customer_to_checking_account"
	"github.com/fwiedmann/ezb/domain/usecase/checking_account_management"
	"github.com/fwiedmann/ezb/domain/usecase/customer_management"
	"github.com/fwiedmann/ezb/domain/usecase/debit"
	"github.com/fwiedmann/ezb/domain/usecase/deposit"
	"github.com/gorilla/mux"
)

// Router is a basic http router with ezb specific REST API
type Router struct {
	CustomerManager          *customer_management.UseCase
	CheckingAccountManager   *checking_account_management.UseCase
	Debit                    *debit.UseCase
	Deposit                  *deposit.UseCase
	Mapper                   *assign_customer_to_checking_account.UseCase
	Port                     int
	GracefullyShutdownTimout time.Duration
}

// StartRouter inits a new HTTP router which implements a REST API for interacting with the ezb service
func (r *Router) StartRouter(stopChan <-chan struct{}, errChan chan<- error) {
	router := mux.NewRouter()

	cuh := &customerHandler{manager: r.CustomerManager}
	router.NewRoute().Path("/customer").HandlerFunc(cuh.CreateCustomer).Methods("POST")
	router.NewRoute().Path("/customer/{id}").HandlerFunc(cuh.UpdateCustomer).Methods("PUT")
	router.NewRoute().Path("/customer/{id}").HandlerFunc(cuh.GetCustomer).Methods("GET")

	cah := checkingAccountHandler{
		manager:        r.CheckingAccountManager,
		depositUseCase: r.Deposit,
		debitUseCase:   r.Debit,
	}
	router.NewRoute().Path("/checking-account").HandlerFunc(cah.CreateCheckingAccount).Methods("POST")
	router.NewRoute().Path("/checking-account/{id}").HandlerFunc(cah.UpdateCheckingAccount).Methods("PUT")
	router.NewRoute().Path("/checking-account/{id}").HandlerFunc(cah.GetCheckingAccount).Methods("GET")
	router.NewRoute().Path("/checking-account/{id}/deposit").HandlerFunc(cah.Deposit).Methods("POST")
	router.NewRoute().Path("/checking-account/{id}/debit").HandlerFunc(cah.Debit).Methods("POST")

	mapperHandler := customerCheckingAccountMapper{mapper: r.Mapper}
	router.NewRoute().Path("/customer-checking-account-mapping").HandlerFunc(mapperHandler.CreateMapping).Methods("POST")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", r.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverErrChan := make(chan error)
	go func() {
		serverErrChan <- server.ListenAndServe()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), r.GracefullyShutdownTimout)
	defer cancel()

	select {
	case err := <-serverErrChan:
		errChan <- err
	case <-stopChan:
		errChan <- server.Shutdown(ctx)
	}
}
