package domain

import "time"

// TelemetryEvent captures metrics for every processed chat completion.
type TelemetryEvent struct {
	ID               string    `json:"id"`
	Timestamp        time.Time `json:"timestamp"`
	Model            string    `json:"model"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	EstimatedCostUSD float64   `json:"estimated_cost_usd"`
	DollarsSavedUSD  float64   `json:"dollars_saved_usd"`
	LatencyMs        int64     `json:"latency_ms"`
	IsCacheHit       bool      `json:"is_cache_hit"`
	Provider         string    `json:"provider"`
	CircuitState     string    `json:"circuit_state"`
}

// CalculateCostEstimates calculates nominal cost and saved cost.
func CalculateCostEstimates(model string, promptTokens, completionTokens int, isCacheHit bool, provider string) (float64, float64) {
	// Standard GPT-4o reference rates ($2.50 / 1M prompt, $10.00 / 1M completion)
	standardCost := (float64(promptTokens) * 0.0000025) + (float64(completionTokens) * 0.0000100)

	if isCacheHit {
		return 0.0, standardCost
	}

	if provider == "gemini-fallback" {
		// Gemini Flash rate (~$0.075 / 1M prompt, $0.30 / 1M completion)
		actualCost := (float64(promptTokens) * 0.000000075) + (float64(completionTokens) * 0.00000030)
		savedCost := standardCost - actualCost
		if savedCost < 0 {
			savedCost = 0
		}
		return actualCost, savedCost
	}

	return standardCost, 0.0
}
