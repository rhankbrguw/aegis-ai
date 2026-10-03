# AegisAI

Resilient AI gateway and real-time SRE mobile sentinel. Go backend reverse proxy with Redis prompt cache, PII guardrails, and Gemini failover routing. Flutter mobile client for real-time telemetry and kill-switch control.

---

Clone it, copy `deployments/.env.example` to `deployments/.env`, then run with Docker or natively.

Backend (Docker): `docker compose -f deployments/docker-compose.yml up -d` (runs gateway on `:8080`, Redis on `:6379`).

Backend (Host): `cd services/gateway && go run ./cmd/server` (requires local Redis).

Client (Flutter): `cd mobile && flutter pub get && flutter run`.

Production Build: `cd mobile && flutter build appbundle --release --dart-define=APP_ENV=production`.

Quality Verification: `make check` (runs `go vet`, `go test -race -cover`, `flutter analyze`, and `flutter test`).

---

Requires Go 1.25+, Flutter 3.49+, Docker, Redis 7.
