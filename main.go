package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fwiedmann/ezb/domain/usecase/assign_customer_to_checking_account"

	"github.com/fwiedmann/ezb/internal/api/http"

	"github.com/fwiedmann/ezb/domain/usecase/checking_account_management"

	"github.com/fwiedmann/ezb/domain/usecase/customer_management"

	"github.com/fwiedmann/ezb/domain/entity/customer_checking_account_mapping"

	"github.com/fwiedmann/ezb/domain/usecase/debit"

	"github.com/fwiedmann/ezb/domain/usecase/deposit"

	"github.com/fwiedmann/ezb/domain/entity/checking_account"

	"github.com/fwiedmann/ezb/domain/entity/customer"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbConnection := os.Getenv("EZB_MYSQL_CONNECTION_STRING")
	if dbConnection == "" {
		log.Fatal("required environment variable \"EZB_MYSQL_CONNECTION_STRING\" is not set or empty")
	}

	db, err := sql.Open("mysql", dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	customerRepo, err := customer.NewMySqlRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	customerManager := customer_management.NewUseCase(customer.NewManager(customerRepo))

	checkinAccountRepo, err := checking_account.NewMySqlRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	checkinAccountRepoManager := checking_account_management.NewUseCase(checking_account.NewManager(checkinAccountRepo))

	depositUseCase := deposit.NewUseCase(checkinAccountRepoManager)
	debitUseCase := debit.NewUseCase(checkinAccountRepoManager)

	mapperRepo, err := customer_checking_account_mapping.NewMySqlRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	mapperUseCase := assign_customer_to_checking_account.NewUseCase(customer_checking_account_mapping.NewManager(mapperRepo))

	r := http.Router{
		CustomerManager:        customerManager,
		CheckingAccountManager: checkinAccountRepoManager,
		Debit:                  debitUseCase,
		Deposit:                depositUseCase,
		Mapper:                 mapperUseCase,
		Port:                   8080,
	}
	stopChan := make(chan struct{})
	httpRouterErrChan := make(chan error)
	go r.StartRouter(stopChan, httpRouterErrChan)
	select {
	case <-initOSNotifyChan():
		stopChan <- struct{}{}
	case err := <-httpRouterErrChan:
		panic(err)
	}
}

func initOSNotifyChan() <-chan os.Signal {
	notifyChan := make(chan os.Signal, 3)
	signal.Notify(notifyChan, syscall.SIGTERM, syscall.SIGINT)
	return notifyChan
}
