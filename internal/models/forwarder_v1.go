package models

type ForwarderV1 struct {
	Version int               `json:"version" default:"1"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}
