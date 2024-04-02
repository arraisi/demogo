include .env

generate-pb-go :
	./scripts/protoc.sh -p $(PROTO)

new-migration:
	migrate create -ext sql -dir scripts/migrations -seq $(name)

migration-up:
	migrate -database "postgres://$(POSTGRES_DATABASE_WRITE_USERNAME):$(POSTGRES_DATABASE_WRITE_PASSWORD)@$(POSTGRES_DATABASE_WRITE_HOST):$(POSTGRES_DATABASE_WRITE_PORT)/$(POSTGRES_DATABASE_WRITE_NAME)?sslmode=disable" -path scripts/migrations up $(step)

migration-down:
	migrate -database "postgres://$(POSTGRES_DATABASE_WRITE_USERNAME):$(POSTGRES_DATABASE_WRITE_PASSWORD)@$(POSTGRES_DATABASE_WRITE_HOST):$(POSTGRES_DATABASE_WRITE_PORT)/$(POSTGRES_DATABASE_WRITE_NAME)?sslmode=disable" -path scripts/migrations down $(step)

migration-force:
	migrate -database "postgres://$(POSTGRES_DATABASE_WRITE_USERNAME):$(POSTGRES_DATABASE_WRITE_PASSWORD)@$(POSTGRES_DATABASE_WRITE_HOST):$(POSTGRES_DATABASE_WRITE_PORT)/$(POSTGRES_DATABASE_WRITE_NAME)?sslmode=disable" -path scripts/migrations force $(version)
