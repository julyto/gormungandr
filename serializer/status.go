package serializer

import (
	"time"
)

type StatusResponse struct {
	Status         string    `json:"status,omitempty"`
	Runtime        string    `json:"runtime,omitempty"`
	Version        string    `json:"version,omitempty"`
	LoadConfAt     time.Time `json:"load_conf_at,omitempty"`
	LoadConfStatus string    `json:"load_conf_status,omitempty"`
}
