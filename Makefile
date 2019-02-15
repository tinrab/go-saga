docker: ## Run containers
	docker-compose -f ./test/docker-compose.yml up -d --build

docker-down: ## Shutdown containers
	docker-compose -f ./test/docker-compose.yml stop
	docker-compose -f ./test/docker-compose.yml rm -f

help: ## Display this help screen
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
