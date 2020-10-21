all: run
ACCOUNT_ID =

run:
	docker-compose up -d --build

down:
	docker-compose down

customer:
	curl -XPOST -d @testUser.json http://localhost:8080/customer

account:
	curl -XPOST -d @testCheckingAccount.json http://localhost:8080/checking-account
debit:
	curl -XPOST -d @testDeposit.json http://localhost:8080/checking-account/$(ACCOUNT_ID)/debit

deposit:
	curl -XPOST -d @testDeposit.json http://localhost:8080/checking-account/$(ACCOUNT_ID)/deposit

mapping:
	curl -XPOST -d @testMapping.json http://localhost:8080/customer-checking-account-mapping