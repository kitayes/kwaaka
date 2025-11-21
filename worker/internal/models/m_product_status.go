package models

import "time"

type ProductStatusRequest struct {
	Status string `json:"status"` // available | not_available | deleted
	Reason string `json:"reason"` // out_of_stock и т.п.
}

type ProductStatusEvent struct {
	EventType string    `json:"event_type"` // "product.status_changed"
	ProductID string    `json:"product_id"`
	OldStatus string    `json:"old_status,omitempty"`
	NewStatus string    `json:"new_status"`
	Reason    string    `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id,omitempty"`
}
