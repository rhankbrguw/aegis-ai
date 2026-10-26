# Task List: AegisAI Enterprise Upgrade

## Phase 1: Real-time Guardrails & PII Redaction (`guardrail-engine`)
- [x] Task 1.1: Define Guardrail Domain Entities & Policy Error Constants
  - Acceptance: Define `GuardrailResult`, `PIIMapping`, and `ErrCodePolicyViolation` constants.
  - Verify: `go test -v ./internal/domain/...`
  - Files: `services/gateway/internal/domain/guardrail.go`, `services/gateway/internal/constants/errors.go`

- [x] Task 1.2: Implement `GuardrailService` for PII Masking and Prompt Injection Detection
  - Acceptance: Detects PII (email, phone, cards, NIK) and jailbreak heuristics; masks sensitive tokens reversibly.
  - Verify: `go test -v ./internal/service/...`
  - Files: `services/gateway/internal/service/guardrail_service.go`, `services/gateway/internal/service/guardrail_service_test.go`

- [x] Task 1.3: Integrate Guardrail Pipeline into `ProxyEngine`
  - Acceptance: Scans incoming requests before dispatch; blocks malicious prompts with HTTP 422 `SECURITY_POLICY_VIOLATION`.
  - Verify: `go test -v ./internal/service/...`
  - Files: `services/gateway/internal/service/proxy_service.go`

---

## Phase 2: SSE Streaming Multiplexer & TTFT Tracker (`streaming-engine`)
- [x] Task 2.1: Implement Server-Sent Events Stream Formatter
  - Acceptance: Formats OpenAI-compliant SSE chunks (`data: {...}\n\n` and `data: [DONE]`).
  - Verify: `go test -v ./internal/domain/...`
  - Files: `services/gateway/internal/domain/streaming.go`

- [x] Task 2.2: Add Streaming Proxy Support to `ChatHandler` & `ProxyEngine`
  - Acceptance: Streams tokens chunk-by-chunk when `stream: true`; emits TTFT latency to telemetry hub.
  - Verify: `go test -v ./internal/handler/...`
  - Files: `services/gateway/internal/handler/chat_handler.go`, `services/gateway/internal/service/streaming_service.go`

---

## Phase 3: Universal Provider Adapter (OpenAI ⇄ Gemini) (`provider-adapters`)
- [x] Task 3.1: Implement Google Gemini `generateContent` Schema Translator
  - Acceptance: Maps OpenAI messages (`system`, `user`, `assistant`) to Gemini `contents` and `parts`.
  - Verify: `go test -v ./internal/service/...`
  - Files: `services/gateway/internal/service/gemini_adapter.go`, `services/gateway/internal/service/gemini_adapter_test.go`

- [x] Task 3.2: Wire Real Gemini Adapter into `FallbackRouter`
  - Acceptance: Executes real HTTP call to Gemini API on primary upstream outage.
  - Verify: `go test -v ./internal/service/...`
  - Files: `services/gateway/internal/service/fallback_router.go`
