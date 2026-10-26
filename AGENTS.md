# AGENTS.md

> Read this file before every task. The global `engineering-standards` skill applies
> to all code. This file provides project-specific context that overrides or extends it.

---

## Project

```
name    : AegisAI (AegisOps)
stack   : Go 1.25 (Microservices) + Flutter 3.49 / Dart 3.13 (Mobile)
arch    : Layered Microservices + Cloudflare Edge + Upstash Redis + TiDB/Neon Serverless
infra   : Docker Compose (Local) + Oracle Cloud Always Free / Cloud Run (Prod)
domain  : rhankbrguw.xyz (api.rhankbrguw.xyz)
```

## Active Stack Rules

```
stacks: [go, flutter, dart, docker]
```

## Folder Structure

```
services/gateway/
├── cmd/server/main.go            # Entrypoint (dependency injection & server startup only)
├── internal/
│   ├── config/                   # Config loader & env validation
│   ├── constants/                # strings.go, routes.go, errors.go, tokens.go
│   ├── domain/                   # Domain entities, interfaces & custom error types
│   ├── handler/                  # HTTP & WebSocket handlers (boundary validation & response envelope)
│   ├── service/                  # Business logic (Circuit Breaker, Proxy engine, Fallback router)
│   └── repository/               # Redis cache & TiDB data access
└── Dockerfile

mobile/
├── lib/
│   ├── main.dart
│   ├── core/
│   │   ├── constants/            # colors.dart, strings.dart, routes.dart, tokens.dart, config.dart
│   │   ├── errors/               # Typed exceptions & failure wrappers
│   │   ├── network/              # Dio HTTP & WebSocket client with interceptors
│   │   └── theme/                # Cyberpunk Dark Theme tokens & design system
│   └── features/
│       ├── dashboard/            # Speedometer, token burn metrics
│       ├── circuit_breaker/      # Kill-switch slider & manual fallback controls
│       └── logs/                 # Realtime streamed logs & payload diff
└── pubspec.yaml
```

## Error Code Registry

| Code | HTTP Status | Meaning |
| :--- | :--- | :--- |
| `VALIDATION_ERROR` | 422 | Input validation failed |
| `UNAUTHENTICATED` | 401 | Missing or invalid API Key / JWT |
| `UNAUTHORIZED` | 403 | Insufficient permissions for operation |
| `NOT_FOUND` | 404 | Resource or route does not exist |
| `CONFLICT` | 409 | Duplicate resource or state conflict |
| `CIRCUIT_TRIPPED` | 503 | Circuit breaker is OPEN (upstream outage / rate limit) |
| `BUDGET_EXCEEDED` | 429 | Daily / Monthly FinOps token limit breached |
| `KILL_SWITCH_ACTIVE` | 503 | Manual emergency kill-switch is active |
| `UPSTREAM_ERROR` | 502 | Upstream LLM provider returned unhandled failure |
| `INTERNAL_ERROR` | 500 | Unexpected internal failure |

## Agent Constraints

### Must:
- Follow `engineering-standards` strictly on every line of code.
- Keep all files **≤ 150 lines** and all functions **≤ 30 lines**.
- Enforce strict 3-tier Layered Architecture: `Handler` -> `Service` -> `Repository`.
- All Go functions doing I/O must accept `ctx context.Context` as the first parameter.
- Use unified constants (`constants/` directory) for all strings, colors, routes, and config keys.
- Wrap all API responses in the standard JSON envelope (`success`, `code`, `message`, `data`, `meta` / `errors`).
- Use `Riverpod` for Flutter state management and `freezed` for immutable models.
- Validate inputs at the handler boundary before calling any service.

### Must Not:
- Write raw inline string literals, magic numbers, or hardcoded URLs/colors.
- Exceed 3 levels of nesting (use early returns and guard clauses).
- Leave any TODO, placeholder text, or AI commentary in final code.
- Write business logic inside `main.go`, HTTP handlers, or Flutter UI widgets.
- Swallow errors silently (`_ = err` or empty `catch`).

