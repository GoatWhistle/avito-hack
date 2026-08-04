# avito-hack

Production-ready starter for a hackathon web application. Everything that is the same in every
project is already built, tested and running: authentication, authorization, database access,
transactions, error contract, pagination, i18n, design tokens, Docker, CI.

The domain part is a **template**, not a product. When the case is announced you replace the
`item` module with your own entity and keep the rest untouched.

---

## Stack

| Layer | Choice |
|---|---|
| Backend | Go 1.23, chi, pgx/v5, goose, slog, JWT (HS256), bcrypt, Prometheus |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Frontend | React 18, TypeScript 5.6 (strict), Vite 5, Ant Design 5, Effector, TanStack Query 5, Axios |
| Forms | React Hook Form + Zod |
| i18n | i18next (ru / en, namespace per feature) |
| Tests | Go testing + testify, Vitest + Testing Library |
| Infra | Docker Compose, Nginx, GitHub Actions, pre-commit |

---

## Quick start

```bash
make init          # creates .env and installs dependencies
# fill in POSTGRES_PASSWORD and JWT_SECRET in .env
make up            # builds and starts the whole stack, waits until healthy
```

- frontend: http://localhost:3000
- backend: http://localhost:8080/api/v1
- health: http://localhost:8080/healthz and http://localhost:8080/readyz
- metrics: http://localhost:8080/metrics

```bash
make lint          # Go linter, ESLint, tsc, file size limit, locale sync
make test          # backend and frontend tests with coverage
make logs SERVICE=backend LOGS_TAIL=50
make down          # stop and remove containers
make clean         # also delete volumes with all data
make migrate       # apply pending migrations
```

---

## What is already implemented

### Backend

**Authentication and authorization**

- `POST /api/v1/auth/register`, `POST /api/v1/auth/login`
- `GET /api/v1/users/me`, `PATCH /api/v1/users/me`
- JWT issuing and parsing, bcrypt password hashing (cost 12)
- Middleware: `Authenticate` (required), `OptionalAuthenticate` (public endpoints with an
  optional actor), `RequireRole`
- Roles: `user`, `moderator`, `admin`

**Infrastructure**

- Config from environment with validation on startup, the process refuses to boot on a bad config
- pgxpool with tuned limits, `TxManager` that carries the transaction through `context`
- Repositories pick the transaction from the context, so use cases never see pgx
- Structured JSON logging with a request id propagated to the client in every error body
- Prometheus metrics (RED per route), `healthz` and `readyz` split
- Graceful shutdown on SIGTERM
- Single error contract, domain errors mapped to HTTP status codes in one place
- Keyset (cursor) pagination, never `OFFSET`
- goose migrations embedded into the binary, applied by a separate one-shot container

**Value objects and helpers**: `Money` (integer kopeks, no floats), `Email`, `password.Hash`,
`Cursor`, `Clock` (injectable time for deterministic tests).

### Frontend

- Feature-Sliced Design with import boundaries enforced by ESLint
- Auth flow: login, registration, profile, session bootstrap from a stored token, automatic
  logout on 401 triggered from the axios interceptor
- Route guards `RequireAuth` and `RequireGuest`, lazy loaded pages
- Typed API layer, query key factory, React Query hooks including infinite pagination
- Effector for client state (filters, theme, session), React Query for server state
- i18n with typed keys, ru/en namespaces, plural rules, `Intl` for dates and currency
- Design tokens as the single source of truth, feeding both the Ant Design theme and CSS
  variables, light and dark themes
- Error codes from the backend are translated on the client, the raw `message` is never shown
- Reusable `FormField`, `ErrorState`, `EmptyState`, `PageSkeleton`

### Quality gates

- `.golangci.yaml` with `depguard` rules that make the layering physical: `domain` cannot import
  pgx, net/http or another module, `app` cannot import HTTP, `api` cannot import infra
- ESLint with `boundaries`, `i18next/no-literal-string`, a ban on hardcoded colors, `max-lines: 250`
- Hard limit of 250 lines per file, checked by `scripts/check-file-length.sh`
- Locale sync check, so a key added to `ru` cannot be forgotten in `en`
- pre-commit: gitleaks, formatters, linters, short tests, conventional commit messages
- GitHub Actions: quality, backend (lint, tests, migrations up and down, govulncheck),
  frontend (types, lint, tests, build, bundle size), stack smoke test through docker compose,
  image publishing to ghcr with a Trivy scan

---

## Conventions

1. **No comments in code.** Not inline, not block, not doc comments. Only tool directives:
   `//go:build`, `//go:generate`, `//nolint:<linter> // reason`, `-- +goose Up/Down`,
   `// eslint-disable-next-line <rule> -- reason`. If code needs an explanation, rename it,
   extract a function or introduce a value object.
2. **The backend is English only.** Identifiers, error strings, log messages, enum values,
   migration names. The backend never returns user-facing text: it returns a stable machine
   `code`, and the frontend renders the translation for the current locale.
3. **250 lines per file, maximum.** One use case per file, one endpoint per handler file.
4. **No literal strings in the UI.** Everything goes through `t()` and JSON locales.
5. **No hardcoded colors, fonts or spacings.** Only tokens from `@/shared/design`.
6. **All routes live under `/api/v1`.** Breaking changes require `v2`, never an edit of `v1`.

---

## Replacing the template once the case is announced

The `item` module is a working reference implementation of a full vertical slice: aggregate with
invariants, status machine, use cases, repository, read model, HTTP layer, tests. Rename it into
your entity instead of writing the structure from scratch.

### Backend, roughly 15 minutes

1. Copy the module directory:

   ```bash
   cp -r src/backend/internal/module/item src/backend/internal/module/order
   ```

2. Replace the identifiers inside the copy: `item` to `order`, `Item` to `Order`.

3. Edit `domain/order.go`: keep the fields your entity actually has and put every business rule
   into the constructor. An invalid aggregate must be impossible to create.

4. Edit `domain/status.go`: describe the real state machine in the `transitions` map. If your
   entity has no lifecycle, delete that file together with `item_behavior.go` and the status
   use case.

5. Update `infra/pg_repository.go` and `infra/mapper.go` to match the new columns.

6. Create a migration:

   ```bash
   cd src/backend && goose -dir migrations create create_orders sql
   ```

   Duplicate the domain invariants as `CHECK` and `FOREIGN KEY` constraints, add an index for
   every query the read model performs, and always write a working `-- +goose Down`.

7. Register the module in `cmd/api/main.go` next to `user` and `item`.

8. Delete the `item` module once nothing references it.

### Frontend, roughly 15 minutes

1. Copy `src/frontend/src/entities/item` into `entities/order` and rename the identifiers.
2. Update `model/types.ts` and `model/mappers.ts` to match the new API contract.
3. Copy the features you need: `item-create`, `item-edit`, `item-status`, `items-filter`.
4. Copy the pages, register them in `src/app/router/routes.tsx` and add the paths to
   `src/shared/config/routes.ts`.
5. Rename the locale namespace `item.json` in both `ru` and `en`, then register it in
   `shared/i18n/resources.ts` and `shared/i18n/i18next.d.ts`.
6. Run `make lint`: the boundary rules and the locale check will catch anything you missed.

### What you almost never need to touch

`shared/` on both sides, the `user` module, middleware, the error contract, pagination, design
tokens, Docker, CI. That part is already finished.

---

## Project layout

```text
avito-hack/
├── .github/workflows/       quality, backend, frontend, stack, docker
├── scripts/                 file length limit, locale sync
├── src/
│   ├── backend/
│   │   ├── cmd/{api,migrate}/
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── module/{user,item}/{domain,app,infra,api}/
│   │   │   ├── shared/{vo,auth,apierr,domainerr,httpx,pagination,middleware,postgres,logger,validate,password,clock}/
│   │   │   └── server/
│   │   └── migrations/
│   └── frontend/
│       └── src/
│           ├── app/{providers,router,styles}/
│           ├── pages/ widgets/ features/ entities/
│           └── shared/{api,config,design,i18n,ui,lib}/
└── tests/                   integration, e2e, load
```

Import direction on the backend: `api -> app -> domain` and `infra -> domain`. A module never
imports the internals of another module, only ports declared in its own `app/ports.go`.

Import direction on the frontend: `app -> pages -> widgets -> features -> entities -> shared`.

---

## Environment variables

Copy `.env.example` to `.env`. `POSTGRES_PASSWORD` and `JWT_SECRET` are required and have no
default. Everything prefixed with `VITE_` is baked into the frontend bundle at build time and is
therefore public, so never put a secret there.
