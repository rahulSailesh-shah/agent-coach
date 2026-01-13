DB_URL	 = ./data/agent-coach.db
DB_DRIVER = sqlite3

dev:
	wails dev

build:
	wails build

migrate-up:
	@if [ ! -f $(DB_URL) ]; then \
		mkdir -p $(dir $(DB_URL)); \
		touch $(DB_URL); \
	fi
	goose $(DB_DRIVER) -dir internal/migrations $(DB_URL) up
	
create-migration:
	goose -dir internal/migrations $(DB_DRIVER) $(DB_URL) create $(name) sql

make clean:
	rm -rf $(DB_URL)