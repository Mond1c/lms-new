package domain

// ProviderRef identifies a VCS provider instance. Kind mirrors the proto
// ProviderKind enum value; domain stays free of generated types.
type ProviderRef struct {
	Kind     int32
	Instance string
}
