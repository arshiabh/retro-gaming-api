
migration:
	@migrate create -ext sql -dir ./cmd/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))
migrate-up:
	migrate -path ./cmd/migrate/migrations -database "postgres://admin:secret@localhost:5555/retro-game?sslmode=disable" up
migrate-down:
	@migrate -path ./cmd/migrate/migrations -database postgres://admin:secret@localhost:5555/retro-game?sslmode=disable down $(filter-out $@,$(MAKECMDGOALS)) 
migrate-version:
	@migrate -path ./cmd/migrate/migrations -database postgres://admin:secret@localhost:5555/retro-game?sslmode=disable version 
