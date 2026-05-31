package domain

import "net/http"

// Temp
type ProviderKind = int32

// Temp
type NormilizedEvent struct{}

type Provider interface {
	Kind() ProviderKind
	Instance() string
	VerifyWebhook(headers http.Header, body []byte) error
	ParseEvent(headers http.Header, body []byte) *NormilizedEvent
	// other methods
}
