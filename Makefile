ifneq (,$(wildcard ./.env))
    include .env
    export
endif


.PHONY: migrate_db
migrate_db:
	flyway -user=${POSTGRES_USER} -password=${POSTGRES_PASSWORD} -url=jdbc:postgresql://localhost:${POSTGRES_LOCAL_PORT}/ -locations=filesystem:sql migrate


.PHONY: clear_db
clear_db:
	echo "This command cleans up the whole database. Press ENTER if you are sure you are not on production, else press Ctrl+C"; \
	read REPLY; \
	docker-compose stop postgres; \
	rm -rf data/postgres || true; \
	docker-compose up -d postgres; \
	until pg_isready -h localhost -p ${POSTGRES_LOCAL_PORT}; do sleep 1; done; \


.PHONY: fresh_db
fresh_db: clear_db migrate_db


.PHONY: up_storages
up_storages:
	docker-compose up -d postgres redis


.PHONY: rebuild
rebuild:
	go mod tidy
	docker builder prune -f
	docker-compose build --force-rm --no-cache


.PHONY: up
up: rebuild
	docker-compose up -d --remove-orphans

.PHONY: protobuf
protobuf:
	protoc -I src/contracts/ src/contracts/*.proto --go_out=src/contracts --go-grpc_out=require_unimplemented_servers=false:src/contracts

.PHONY: docs
docs:
	cd src/services/api && swag init --parseDependency --parseInternal
