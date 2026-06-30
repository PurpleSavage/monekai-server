# Migraciones en Go con `golang-migrate`

## ¿Qué es una migración?

Una migración es un cambio controlado y versionado en el esquema de la base de datos. En vez de modificar la BD a mano, cada cambio se guarda como un archivo SQL numerado que puede aplicarse o revertirse de forma reproducible en cualquier entorno (local, staging, producción).

---

## Archivos `up` y `down`

Cada migración se compone de **dos archivos**:

| Archivo | Propósito |
|---|---|
| `*.up.sql` | Aplica el cambio — crea tablas, agrega columnas, crea índices |
| `*.down.sql` | Revierte el cambio — hace el opuesto exacto del `up` |

### Ejemplo

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR NOT NULL
);
```

```sql
-- 000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

La regla es simple: **el `down` debe deshacer exactamente lo que hace el `up`**. Si el `up` crea una tabla, el `down` la dropea. Si el `up` agrega una columna, el `down` la elimina.

---

## Nomenclatura

El formato estándar es:

```
{version}_{descripcion}.up.sql
{version}_{descripcion}.down.sql
```

- **version** — número secuencial con ceros a la izquierda (`000001`, `000002`, `000003`)
- **descripcion** — snake_case describiendo qué hace la migración

### Ejemplos reales

```
000001_create_users_table.up.sql
000001_create_users_table.down.sql

000002_create_sessions_table.up.sql
000002_create_sessions_table.down.sql

000003_add_credits_to_users.up.sql
000003_add_credits_to_users.down.sql

000004_create_samples_table.up.sql
000004_create_samples_table.down.sql
```

> El número debe ser único y secuencial. `golang-migrate` los ejecuta en orden ascendente.

---

## ¿Cuándo crear una nueva migración?

Crear una nueva migración cada vez que necesites:

- Crear una tabla nueva
- Eliminar una tabla
- Agregar una columna
- Eliminar una columna
- Modificar el tipo de una columna
- Crear o eliminar un índice
- Agregar o quitar constraints (FK, UNIQUE, CHECK)
- Insertar datos semilla (seed data)

**Nunca edites una migración que ya corrió en producción.** Si necesitas cambiar algo, crea una nueva migración que aplique la corrección.

---

## Cómo funciona internamente

`golang-migrate` mantiene una tabla `schema_migrations` en tu base de datos:

```
schema_migrations
┌─────────┬───────┐
│ version │ dirty │
├─────────┼───────┤
│       1 │ false │
│       2 │ false │
└─────────┴───────┘
```

- **version** — el número de la última migración ejecutada
- **dirty** — `true` si la migración falló a mitad, `false` si está limpia

Cuando corres `m.Up()`, `golang-migrate` revisa esta tabla y solo ejecuta las migraciones cuyo número sea mayor al último registrado. Por eso nunca corre la misma migración dos veces.

---

## Comandos principales

### En código Go

```go
// Correr todas las migraciones pendientes
m.Up()

// Revertir todas las migraciones
m.Down()

// Avanzar N migraciones
m.Steps(2)

// Retroceder N migraciones
m.Steps(-2)

// Ir a una versión específica
m.Migrate(3)
```

### Con el CLI (opcional)

```bash
# Instalar el CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Crear una nueva migración
migrate create -ext sql -dir migrations -seq add_photo_to_users

# Aplicar migraciones
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

# Revertir la última migración
migrate -path migrations -database "postgres://..." down 1

# Ver versión actual
migrate -path migrations -database "postgres://..." version
```

---

## Estado `dirty`

Si una migración falla a mitad de ejecución, `golang-migrate` marca `dirty = true`. En ese caso la app no arrancará hasta que resuelvas el conflicto manualmente:

```bash
# Forzar la versión a un número específico para limpiar el estado dirty
migrate -path migrations -database "postgres://..." force 1
```

Después de forzar, corriges el SQL y vuelves a correr `up`.

---

## Estructura recomendada en el proyecto

```
moneka/
├── migrations/
│   ├── 000001_create_tables.up.sql
│   ├── 000001_create_tables.down.sql
│   ├── 000002_add_photo_url_to_users.up.sql
│   ├── 000002_add_photo_url_to_users.down.sql
│   └── ...
├── cmd/
│   └── scripts/
│       └── migrations.go
├── internal/
└── docker-compose.yml
```

---

## Imports necesarios en Go

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres" // driver de postgres
    _ "github.com/golang-migrate/migrate/v4/source/file"       // driver de archivos locales
)
```

Los imports con `_` son de efecto secundario — registran los drivers internamente sin exponer ninguna función. Sin ellos `golang-migrate` no sabe leer archivos `.sql` ni conectarse a postgres.
