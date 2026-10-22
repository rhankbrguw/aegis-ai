# Implementation Plan: Industry-Grade Enterprise AI Gateway Upgrade

## Overview
Elevate AegisAI to tier-1 enterprise standards by implementing three mission-critical capabilities:
1. **Zero-Latency In-line Guardrails & PII Redactor**: Edge PII masking/unmasking (GDPR/Data Privacy) and Prompt Injection defense.
2. **SSE Streaming Engine (`stream: true`)**: Real-time token streaming passthrough with TTFT (Time-to-First-Token) telemetry.
3. **Universal LLM Schema Adapter**: Native OpenAI ⇄ Google Gemini `generateContent` translation for authentic multi-provider failover.

---

## Architecture Decisions
- **Guardrail Pipeline:** Executed at the Gateway service layer before prompt caching and before upstream dispatch. Replaces sensitive PII with tokens (`[REDACTED_EMAIL_1]`) and blocks known injection attacks with `SECURITY_POLICY_VIOLATION` (422).
- **SSE Stream Protocol:** For `stream: true`, handler sets `text/event-stream` and streams chunked delta payloads (`data: {...}\n\n`) directly with `http.Flusher`. Measures TTFT on first chunk emission.
- **Provider Transformer:** Decouples OpenAI message structures from Gemini `contents.parts` structures, mapping system instructions and user roles cleanly.
- **Constraint Compliance:** Every file strictly ≤ 150 lines, all functions ≤ 30 lines, zero hardcoded tokens.

---

## Task List

### Phase 1: Real-time Guardrails & PII Redaction (`guardrail-engine`)
- [ ] Task 1.1: Define Guardrail Entities, Regex Patterns & Policy Error Codes
- [ ] Task 1.2: Implement `GuardrailService` (PII Masking/Unmasking & Prompt Injection Detection)
- [ ] Task 1.3: Wire Guardrails into `ProxyEngine` & Add Comprehensive Unit Tests

### Checkpoint 1: Guardrails Verification
- [ ] `go test -v -race ./internal/service/...` passes 100% with injection and PII test cases

### Phase 2: SSE Streaming Multiplexer & TTFT Tracker (`streaming-engine`)
- [ ] Task 2.1: Implement Server-Sent Events Chunk Processor & Stream Formatter
- [ ] Task 2.2: Implement Streaming Proxy Route in `ChatHandler` & TTFT Telemetry Emission

### Checkpoint 2: Streaming Verification
- [ ] `go test -v -race ./internal/handler/...` validates chunked SSE response format

### Phase 3: Universal Provider Adapter (OpenAI ⇄ Gemini) (`provider-adapters`)
- [ ] Task 3.1: Implement Gemini `generateContent` Request/Response Schema Mapper
- [ ] Task 3.2: Integrate Gemini Native Fallback Client into `FallbackRouter`

### Checkpoint 3: Complete Quality Gate
- [ ] `make check` passes 100% across Go backend and Flutter mobile client

---

## Risks and Mitigations
| Risk | Impact | Mitigation |
| :--- | :--- | :--- |
| Regex overhead on large prompts | Low | Pre-compile regex matchers at init time, cap prompt scan length |
| Streaming disconnects mid-response | Med | Detect client context cancellation (`ctx.Done()`) and abort upstream immediately |
| Discrepancy between Gemini & OpenAI role schema | Low | Map `assistant` -> `model` and `system` -> `systemInstruction` explicitly |
