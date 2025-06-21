include .env
export

# Configuración
DB_URL ?= $(DB_DSN)
MIGRATIONS_DIR := ./migrate/migrations


# Comandos

# Ejecuta migraciones hacia arriba
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

# Ejecuta migraciones hacia abajo (rollback)
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

# Crear una nueva migración (vacía)
migration:
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "❌ Debes proporcionar un nombre: make migration MIGRATION_NAME=nombre_migracion"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create $(MIGRATION_NAME) sql

# Ver el estado de las migraciones
status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

# Ejecutar hasta una versión específica
migrate-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "Necesitas especificar VERSION=x"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" $(VERSION)

# Limpiar migraciones (opcional, si usas archivos temporales)
clean:
	rm -f *.log

.PHONY: migrate-up migrate-down migration status migrate-to clean

