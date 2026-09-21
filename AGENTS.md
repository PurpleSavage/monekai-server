# Monekai — Contexto del proyecto

Monorepo de dos proyectos vecinos que forman **Monekai**, una plataforma web para
generar y editar samples de audio con IA, compartirlos en comunidad, y gestionar
créditos con pagos.

## Ubicación de los proyectos

| Proyecto | Ruta | Stack |
|---|---|---|
| Backend | `C:\personal-projects\monekai-backend` (workspace actual) | Go 1.25 + Chi + GORM + PostgreSQL |
| Frontend | `C:\personal-projects\monekai-frontend` | Angular 21 (standalone) + Signals + Tailwind v4 |

Ambos proyectos se trabajan desde el workspace del backend. Al trabajar con
frontend, leer/modificar archivos bajo `C:\personal-projects\monekai-frontend`.

---

## Backend (`monekai-backend`)

### Arreglo y librerías
- Go **1.25.8**, router **chi/v5**, **GORM** + driver PostgreSQL (la única BD).
- Auth: **Google OAuth** (idtoken) + **JWT** (`golang-jwt/jwt/v5`).
- **Replicate** (generación audio IA, modelo `meta/musicgen`), **Paddle** (pagos,
  Sandbox), **Cloudflare R2** (storage, S3-compatible, URLs presignadas).
- **SSE** para notificaciones en tiempo real; webhooks verificados (svix/Paddle).
- Docker: `docker-compose.yml` con `db` (postgres:15-alpine), `app`, `ngrok`.

### Estructura
```
cmd/entrypoint/main.go          # bootstrap: router Chi, middlewares, monta módulos
configurations/migrations/      # migraciones SQL (3 versiones up/down)
configurations/persistence/     # modelos GORM
modules/{community,notifications,payments,sampler}/  + shared (auth, common)
```

### Arquitectura
**Hexagonal (Ports & Adapters) + DDD** (reglas en `ROO_CODE_RULES.md`). Cada
módulo sigue la misma estructura:
- `application/` → `dtos/`, `ports/`, `usecases/`
- `domain/` → `entities/`, `valueobjects/`, `enums/`
- `infrastructure/` → `controllers/`, `middlewares/`, `in-adapters/`,
  `out-adapters/` (repos), `mappers/`, `raws/`

Reglas clave: constructores devuelven interfaces, módulos solo dependen de
`shared`, comunicación inter-módulos vía Observer (`ObserverBucket`),
middlewares shared se inicializan en el entrypoint.

### Endpoints (base `http://localhost:8080`)
- **Auth** `/auth`: `POST /login` (token Google), `GET /refresh-token`,
  `GET /profile` (cookie `session_token`; access JWT expira en 1m).
- **Audio** `/audio`: `POST /create` (requiere 20 créditos, asíncrono vía Replicate),
  `GET /samples`, `GET /edit-samples[/{id}]`, `POST /share-sample`,
  `PATCH /edit-samples/{id}/url|effects`, `POST /webhook/songs`, `GET /presigned-url`.
- **Community** `/community`: `GET /samples`, `GET /edit-samples`,
  `GET /latest-samples`, `GET /latest-edit-samples` (cache TTL 3000s),
  `PATCH /like/{sampleID}`, `GET /download/{sampleID}`.
- **Payments** `/payments`: `POST /create?packageId=`, `POST /webhook`.
- **Notifications** `/notifications`: `GET /all`, `PATCH /read-all`, `PATCH /{id}/read`.
- **SSE** `/sse/stream`: eventos `sample_ready`, `sample_error`, `payment_success`,
  `payment_failed` (auth por cookie).

### Notas
- Sin tests en el repo.
- Cuidado con typos en nombres de archivos/carpetas al buscar (ej. `out-dapters`,
  `reponses`, `throtler_request_middleware.go`).
- Paddle en Sandbox; `.env` contiene secrets de dev (no commitear).

---

## Frontend (`monekai-frontend`)

### Stack
- **Angular 21** (standalone, sin NgModules, **zoneless**), **Signals** +
  `rxjs-interop`, **Vitest** + jsdom para tests.
- UI: Tailwind v4 (tema en `styles.css`: `--color-surface`, `--color-brand-primary`),
  `@lucide/angular`, `@taiga-ui/layout`, `ngx-sonner`, efectos `@omnedia`.
- Audio/3D: `wavesurfer.js`, `three`, Web Audio API (`AudioEffectsEngineService`).
- Persistencia local: **Dexie** (IndexedDB, tablas `samplesEdited`, `metadata`).
- Estado/HTTP: servicios que implementan puertos; interceptores funcionales.

### Estructura
```
src/app/
  landing/                      # página de marketing pública
  app.routes.ts                 # rutas: `/`, `/auth/login`, `/monekai/...` (protegidas)
  core/                         # arquitectura hexagonal
    framewrok-utilities/guards|interceptors   # [ojo: typo "framewrok"]
    shared/auth/                # dominio auth (application/domain/infrastructure/state-manager/ui)
    shared/common/              # infra común (errores HTTP, IndexedDB, storage, UI compartida)
    sampler/                    # dominio principal: generación/edición audio
    payments/                   # créditos / billing (Paddle)
    account/                    # cuenta
    aggregates/{community,notifications}/   # for-you + notificaciones SSE
```

Cada módulo de `core/` sigue: `application/` (use-cases, ports, dtos),
`domain/` (entities, value-objects), `infrastructure/` (http, mappers,
persistence), `state-manager/` (signals), `ui/` (pages/components).

### Conexión con el backend
- URL base `http://localhost:8080` en `src/environments/environment.development.ts`
  (el interceptor `apiBaseUrlInterceptor` la antepone y activa `withCredentials`).
- Auth: Google One Tap → `POST /auth/login`; sesión en `AuthStateManager` +
  `localStorage['user-data']`; interceptor de refresh con *single-flight*
  (`BehaviorSubject`) → `GET /auth/refresh-token`; guard valida con `/profile`.
- Clients HTTP que implementan puertos: `AuthPort`, `SamplerPort`,
  `CommunityPort`, etc. (registrados en `app.config.ts`).

### Rutas principales
`/` (landing), `/auth/login`, `/monekai/{sampler,for-you,billing,account}`,
todas las de `/monekai` protegidas por `authGuard` y con lazy loading.

### Notas
- Solo hay un test (`src/app/app.spec.ts`).
- Interceptores y `auth-page` importan `environment.development` directamente
  (el build de producción usaría `http://localhost:8080` — pendiente de corregir).
- `googleClientId` hardcodeado en `environment.development.ts`.
- Commits en español (respetar estilo).
- Typos conocidos en nombres: `framewrok-utilities`, `settins-mode.options.ts`,
  `prompt-imput`, `colections-sample-gallery`, `create-sample-requet.dto.ts`.

---

## Comandos útiles

### Backend
- `docker compose up -d db` — levantar PostgreSQL.
- `go run cmd/entrypoint/main.go` — correr servidor (puerto 8080).
- Migraciones: `cmd/scripts/run_migrations.go` (ver `docs/MIGRATIONS.md`).

### Frontend
- `pnpm install` — instalar dependencias.
- `pnpm start` (o `ng serve`) — dev server en `http://localhost:4200/`.
- `pnpm test` (o `ng test`) — tests con Vitest.
- `pnpm build` (o `ng build`) — build de producción.