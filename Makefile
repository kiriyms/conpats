.PHONY: test
test: ## Run tests
	go test -race -v ./...

.PHONY: race-test
race-test:
	go test -race -v -failfast -count=100 ./...