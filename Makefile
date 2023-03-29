#include .env if exists
-include .env

SUBSCRIPTION_BINARY=subscription-service
PSQL_CONN=host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASS) dbname=$(DB_NAME) port=${DB_PORT} sslmode=disable 

tools: ## Install general tools globally (not in the project)
	go get -u github.com/grailbio-external/goose/cmd/goose
## go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

# up_build will run docker-compose building and then running process
up_build: 
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

migrateup: 
#migrate -path resources/database/migration/ -database "postgresql://postgres:password@localhost:5433/subscription?sslmode=disable" -verbose up
	goose -dir resources/database/migration/ postgres "${PSQL_CONN}" up

migratedown:
	goose -dir resources/database/migration/ postgres "${PSQL_CONN}" down

.PHONY: migratedown migrateup

generate_data:
	bash generate_data.sh

generate_one_sub_post:
	bash generate_single_request.sh

clean-cache:
	go clean -testcache


copy_data:
## first, need to copy the generated files to the Postgres docker container
	docker cp GenerateData/Generated/. subscription-postgres:/myData
	docker cp copy_data_to_postgres.sql subscription-postgres:/

## second, execute the sql file in the Postgres docker container
	docker exec -it subscription-postgres psql -U postgres -q -f /copy_data_to_postgres.sql


# this one is only used to build the app, if you use make up_build, this command does not need to run
build_subscription:
	@echo "Building playlist binary..."
	cd ../subscription-service && env GOOS=linux CGO_ENABLED=0 go build -o ${SUBSCRIPTION_BINARY} ./cmd
	@echo "Done!"

container-bridge:
	bash create_docker_bridge.sh

rm-bridge:
	docker network rm multi-host-network
