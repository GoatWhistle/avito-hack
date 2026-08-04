# TASK Plan: Backend #3 (Controllers, Transport Adapters & Integration)

This task list tracks the progress and detailed roadmap for **Backend #3** (Raccoon Gamification module HTTP endpoints, adapters, and event handlers).

---

## ✅ Completed Tasks

- [x] **Architecture Cleanup & Alignment**
  - [x] Removed cloned item CRUD files from `raccoon` module (`domain/`, `infra/`, `write_handlers.go`).
  - [x] Defined clean interfaces in `app/ports.go` (`RaccoonStateReader`, `CoreRaccoonService`, `PromocodeGeneratorAdapter`, `NotificationAdapter`).
  - [x] Implemented development stubs in `app/stub_core.go`.

- [x] **Profile Endpoint (`GET /api/v1/raccoon/profile`)**
  - [x] Implemented `GetProfile` handler invoking `GetRaccoonProfileUseCase`.
  - [x] Mapped `RaccoonProfileView` to `RaccoonProfileResponse` DTO (level, XP, streak, badges).
  - [x] Registered route with `Authenticate` middleware in `api/routes.go`.

- [x] **Claim Reward Endpoint (`POST /api/v1/rewards/claim`)**
  - [x] Created `ClaimRewardRequest` & `ClaimRewardResponse` DTOs.
  - [x] Implemented `ClaimReward` handler invoking `ClaimRewardUseCase`.
  - [x] Registered route with `Authenticate` middleware in `api/routes.go`.

- [x] **Badges Endpoint (`GET /api/v1/badges`)**
  - [x] Implemented `GetBadges` handler invoking `CoreRaccoonService.EvaluateBadges`.
  - [x] Returned mapped `BadgeResponse` array.
  - [x] Registered route with `Authenticate` middleware in `api/routes.go`.

---

## 🚀 Next Steps & Detailed Implementation Plan

### 1. User Action Endpoint (`POST /api/v1/raccoon/action`)
- [ ] **DTOs (`api/dto.go`)**:
  - `UserActionRequest`: `action_type` (string, required), `payload` (map[string]any).
  - `ProcessActionResultResponse`: `xp_added`, `current_xp`, `current_level`, `level_up`, `current_streak`, `new_badges`.
- [ ] **UseCase / Handler (`api/read_handlers.go` or `api/action_handler.go`)**:
  - Extract `Actor` from request context.
  - Call `deps.Core.CheckRateLimit(ctx, actor.ID, req.ActionType)`. If limited, return HTTP 429 / RateLimit error.
  - Call `deps.Core.ProcessUserAction(ctx, ActionEvent{...})`.
  - Return `ProcessActionResultResponse` with HTTP 200 OK.
- [ ] **Route Registration (`api/routes.go`)**:
  - Add `POST /action` under `r.Route("/raccoon", ...)` protected by `h.deps.Authenticate`.

### 2. Infrastructure Adapters (`infra/`)
- [ ] **Promocode Generator (`infra/promocode_adapter.go`)**:
  - Replace stub with formal promocode generator adapter implementing `PromocodeGeneratorAdapter`.
  - Support reward types (e.g. `AVITO-XL-[REWARD]-[UUID]`, delivery discounts, badge perks).
- [ ] **In-App Push / Notification Adapter (`infra/notification_adapter.go`)**:
  - Implement `NotificationAdapter` interface to log or store in-app notifications when rewards or badges are unlocked.

### 3. Asynchronous Event Consumer (`transport/consumer/`)
- [ ] **Event Listener for System Actions**:
  - Implement consumer handler for user platform events (e.g. item listed, deal completed, review written).
  - Parse event payload and invoke `CoreRaccoonService.ProcessUserAction` asynchronously.

### 4. HTTP API Unit & Integration Tests (`api/handlers_test.go`)
- [ ] Write unit tests using `net/http/httptest` for:
  - `GET /api/v1/raccoon/profile` (200 OK, 401 Unauthorized)
  - `GET /api/v1/badges` (200 OK, 401 Unauthorized)
  - `POST /api/v1/rewards/claim` (200 OK, 400 Bad Request, 401 Unauthorized)
  - `POST /api/v1/raccoon/action` (200 OK, 429 Too Many Requests)

### 5. Final Integration with Backend #2
- [ ] When Backend #2 completes Postgres database models for Raccoon state:
  - Replace `StubRaccoonStateReader` and `StubCoreRaccoonService` in `module.go` with real Postgres-backed implementations.
  