package models

type ForwarderV1 struct {
	Version int               `json:"version"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
