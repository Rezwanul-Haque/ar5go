PROJECT_NAME := ar5go
PKG_LIST := $(shell go list ${PROJECT_NAME}/... | grep -v /vendor/)


.PHONY: all dep build ar5go test

all: build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


########################
### DEVELOP and TEST ###
########################
install-swagger:
	which swagger || go install github.com/go-swagger/go-swagger/cmd/swagger

swagger: install-swagger
	swagger generate spec -o ./swagger.yaml --scan-models

development: swagger
	# booting up dependency containers
	@docker-compose up -d consul db redis

	# wait for consul container be ready
	@while ! curl --request GET -sL --url 'http://localhost:8500/' > /dev/null 2>&1; do printf .; sleep 1; done

	# setting KV, dependency of app
	@curl --request PUT --data-binary @config.local.json http://localhost:8500/v1/kv/${PROJECT_NAME}

	# building ar5go
	@docker-compose up --build ${PROJECT_NAME}

test: ## Run unittests
	@go test -cover -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@go tool cover -func=cov.out

ar5go: ## Remove previous build
	@rm -f $(PROJECT_NAME)
	@docker-compose down

