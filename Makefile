include .env
export

# Configuración
DB_URL ?= $(DB_DSN)
MIGRATIONS_DIR := ./migrate/migrations


# Comandos

# Crear una nueva migración (vacía)
migration:
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "❌ Debes proporcionar un nombre: make migration MIGRATION_NAME=nombre_migracion"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create $(MIGRATION_NAME) sql

# Ejecuta migraciones hacia arriba
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

# Ejecutar hasta una versión específica
migrate-up-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "Necesitas especificar VERSION=x"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" $(VERSION)

# Ejecuta solo una migración hacia arriba (up by one)
migrate-up-by-one:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up-by-one

# Ejecuta migraciones hacia abajo (rollback)
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

# Revertir migraciones hasta una versión específica (down to)
migrate-down-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "Necesitas especificar VERSION=x"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down-to $(VERSION)

# Ver el estado de las migraciones
status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

# Limpiar migraciones (opcional, si usas archivos temporales)
clean:
	rm -f *.log

.PHONY: migrate-up migrate-down migration status migrate-to clean

