package serializer

import (
	"time"
)

type StatusResponse struct {
	Status     string    `json:"status,omitempty"`
	Version    string    `json:"version,omitempty"`
	Runtime    string    `json:"runtime,omitempty"`
	LoadConfAt time.Time `json:"load_conf_at,omitempty"`
}
