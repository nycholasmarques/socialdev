# Run sqlc to generate the queries.go file.
sqlc:
	sqlc generate

# Run atlas to migrate the database based on the schema.sql file.
atlas:
	atlas schema apply --env development

lint:
	go fmt ./...