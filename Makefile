MIGRATIONS_DIR := ./migrations

.PHONY: migrate

migrate:
	cd $(MIGRATIONS_DIR) && tern migrate
