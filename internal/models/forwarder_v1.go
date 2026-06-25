package models

// ForwarderV1 is the legacy (version 1) forwarder shape: a bare HTTP URL with
// headers. It is retained only so old forwarders can be read and upgraded to V2
// on load (see forwarder_convert.go).
type ForwarderV1 struct {
	Version int               `json:"version" default:"1"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
