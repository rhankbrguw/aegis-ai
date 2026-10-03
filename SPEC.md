# Spec: AegisAI (AegisOps) — Resilient AI Gateway & Mobile Sentinel

## 1. Capability Map

| Module ID | Responsibility | Depends On |
| :--- | :--- | :--- |
| `gateway-core` | OpenAI-compatible reverse proxy, Circuit Breaker state machine, multi-LLM dynamic failover | — |
| `cache-telemetry` | Redis semantic prompt cache, token counter, async telemetry event stream | `gateway-core` |
| `finops-worker` | Asynchronous cost aggregation, budget threshold checker, anomaly detection | `cache-telemetry` |
| `mobile-sentinel` | Flutter mobile dashboard, live WebSocket telemetry stream, biometric kill-switch | `gateway-core`, `finops-worker` |

**Build Order:** `gateway-core` → `cache-telemetry` → `finops-worker` → `mobile-sentinel`

---

## 2. Assumptions
1. **Primary Upstream:** OpenAI API (`https://api.openai.com/v1`) using standard `/v1/chat/completions` spec.
2. **Fallback Upstream:** Google Gemini 2.5 Flash (`https://generativelanguage.googleapis.com/v1beta`) mapped to OpenAI schema.
3. **Storage & Cache:** Redis 7+ (Local Docker / Upstash Serverless) for caching, rate limits, and circuit state.
4. **Mobile Target:** Android (SDK 34+), responsive cross-platform Flutter client with Biometric support.
5. **Network / Domain:** Traffic routed via Cloudflare Edge to `api.rhankbrguw.xyz`.

---

## 3. Objective & Success Criteria
- **Objective:** Provide zero-code drop-in resilience for AI applications with sub-15ms semantic caching, automatic failover when OpenAI returns 429/5xx, and real-time mobile command controls.
- **Success Criteria:**
  1. `POST /v1/chat/completions` responds in `< 15ms` on Redis cache hit with cost `$0`.
  2. Circuit breaker trips from `CLOSED` to `OPEN` after 5 consecutive failures, routing traffic to Gemini without dropping user requests.
  3. Mobile emergency kill-switch flips proxy state in `< 200ms` via WebSockets.
  4. 100% compliance with `engineering-standards` (all files ≤ 150 lines, functions ≤ 30 lines, zero hardcoded tokens).

---

## 4. Tech Stack & Commands
- **Backend:** Go 1.25, Chi Router, go-redis v9, go-playground/validator v10.
- **Mobile:** Flutter 3.49 / Dart 3.13, Riverpod 2.x, Dio 5.x, Freezed.
- **Infrastructure:** Docker Compose (local), Cloudflare DNS/SSL (`rhankbrguw.xyz`).

```bash
# Gateway Commands
Test:   cd services/gateway && go test -v -race -cover ./...
Lint:   cd services/gateway && golangci-lint run
Run:    docker compose -f deployments/docker-compose.yml up --build

# Mobile Commands
Test:   cd mobile && flutter test
Run:    cd mobile && flutter run
Build:  cd mobile && flutter build appbundle --release
```

---

## 5. Architecture & Code Conventions

### Go Layered Architecture (Backend)
```
Request → internal/handler (validate & envelope)
               ↓
          internal/service (circuit breaker, proxy logic, failover)
               ↓
          internal/repository (redis cache, event stream)
```

### Flutter Clean Architecture (Mobile)
```
UI Widget → Presentation (Riverpod StateNotifier)
                 ↓
            Domain (Entities & UseCases)
                 ↓
            Data (Dio/WebSocket Repository & Local Storage)
```

---

## 6. Boundaries

- **Always Do:**
  - Wrap all API responses in standard JSON envelope (`success`, `code`, `message`, `data`, `meta`).
  - Pass `ctx context.Context` as first parameter in all Go I/O methods.
  - Externalize all strings, colors, routes, and configs into `constants/`.
- **Ask First:**
  - Adding third-party dependencies outside approved list.
  - Modifying the public `/v1/chat/completions` API schema.
- **Never Do:**
  - Hardcode API keys, secrets, or raw color hex values.
  - Exceed 150 lines per file or 30 lines per function.
  - Swallow errors silently.
