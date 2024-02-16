DATABASE_URL = postgres://ranmerc@localhost:5432/moviepin?sslmode=disable
MIGRATIONS_PATH = db/migrations

create-migration:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) create_db

migrate-up:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) up

migrate-down:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) -verbose down
