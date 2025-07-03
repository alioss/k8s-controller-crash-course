package models

import (
	"time"
)

// AtlasApp represents an Atlas application deployment
type AtlasApp struct {
	Namespace   string    `json:"namespace"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	MigrationID string    `json:"migration_id"`
	Status      string    `json:"status"`
	LastUpdate  time.Time `json:"last_update"`
	Age         string    `json:"age"`
	Replicas    string    `json:"replicas"` // e.g., "3/3"
}

// AtlasAppList represents a list of Atlas applications
type AtlasAppList []AtlasApp

// StatusType represents different deployment statuses
type StatusType string

const (
	StatusRunning  StatusType = "Running"
	StatusPending  StatusType = "Pending"
	StatusFailed   StatusType = "Failed"
	StatusUnknown  StatusType = "Unknown"
)

// GetStatus returns a formatted status based on ready replicas
func GetStatus(ready, total int32) StatusType {
	if ready == total && ready > 0 {
		return StatusRunning
	} else if ready == 0 {
		return StatusFailed
	} else if ready < total {
		return StatusPending
	}
	return StatusUnknown
}
