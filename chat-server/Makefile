.PHONY: migrate-up migrate-down

migrate-up:
	@migrate -database "postgres://tester:123@localhost:5432/hw6?sslmode=disable" -path ./migrations up
	@echo "Migrations applied successfully."

migrate-down:
	@migrate -database "postgres://tester:123@localhost:5432/hw6?sslmode=disable" -path ./migrations down
	@echo "Migrations rolled back successfully."
