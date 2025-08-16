# Run sqlc to generate the queries.go file.
sqlc:
	sqlc generate

# Run atlas to migrate the database based on the schema.sql file.
atlas:
	atlas migrate apply --env local

# Generate a new Atlas migration
diff:
	@if [ -z "$(NAME)" ]; then echo "Usage: make diff NAME=<migration_name>"; exit 1; fi
	atlas migrate diff $(NAME) --env local --to file://internal/db/schema.sql --dir file://internal/db/migrations

lint:
	go fmt ./...