package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// GuardrailService manages prompt injection detection and PII redaction.
type GuardrailService struct {
	emailRegex     *regexp.Regexp
	phoneRegex     *regexp.Regexp
	cardRegex      *regexp.Regexp
	nikRegex       *regexp.Regexp
	keyRegex       *regexp.Regexp
	injectionRegex *regexp.Regexp
}

// NewGuardrailService initializes pre-compiled security regex patterns.
func NewGuardrailService() *GuardrailService {
	return &GuardrailService{
		emailRegex:     regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`),
		phoneRegex:     regexp.MustCompile(`(?:\+?62|08)[0-9]{8,12}\b`),
		cardRegex:      regexp.MustCompile(`\b(?:\d{4}[-\s]?){3}\d{4}\b`),
		nikRegex:       regexp.MustCompile(`\b[1-9][0-9]{15}\b`),
		keyRegex:       regexp.MustCompile(`(?:sk-[a-zA-Z0-9]{20,}|ghp_[a-zA-Z0-9]{20,}|aegis_live_[a-zA-Z0-9]{20,})`),
		injectionRegex: regexp.MustCompile(`(?i)(ignore\s+(all\s+)?previous\s+instructions|system\s+override\s+mode|dan\s+mode\s+enabled|disregard\s+(the\s+)?system\s+prompt|bypass\s+all\s+safety)`),
	}
}

// InspectAndSanitize scans messages for injection and masks sensitive PII.
func (s *GuardrailService) InspectAndSanitize(messages []domain.ChatMessage) *domain.GuardrailResult {
	redactions := make(map[string]string)
	sanitized := make([]domain.ChatMessage, len(messages))

	for i, msg := range messages {
		if s.injectionRegex.MatchString(msg.Content) {
			return &domain.GuardrailResult{
				SanitizedMessages: messages,
				IsBlocked:         true,
				BlockReason:       "Prompt injection / jailbreak pattern detected",
			}
		}
		sanitizedContent := s.maskPII(msg.Content, redactions)
		sanitized[i] = domain.ChatMessage{Role: msg.Role, Content: sanitizedContent}
	}

	return &domain.GuardrailResult{
		SanitizedMessages: sanitized,
		Redactions:        redactions,
		IsBlocked:         false,
	}
}

// UnmaskResponse restores redacted sensitive placeholders in output.
func (s *GuardrailService) UnmaskResponse(content string, redactions map[string]string) string {
	for placeholder, original := range redactions {
		content = strings.ReplaceAll(content, placeholder, original)
	}
	return content
}

func (s *GuardrailService) maskPII(text string, store map[string]string) string {
	text = s.replacePattern(text, s.emailRegex, "EMAIL", store)
	text = s.replacePattern(text, s.phoneRegex, "PHONE", store)
	text = s.replacePattern(text, s.cardRegex, "CARD", store)
	text = s.replacePattern(text, s.nikRegex, "NIK", store)
	return s.replacePattern(text, s.keyRegex, "KEY", store)
}

func (s *GuardrailService) replacePattern(text string, re *regexp.Regexp, tag string, store map[string]string) string {
	return re.ReplaceAllStringFunc(text, func(match string) string {
		placeholder := fmt.Sprintf("[REDACTED_%s_%d]", tag, len(store)+1)
		store[placeholder] = match
		return placeholder
	})
}
