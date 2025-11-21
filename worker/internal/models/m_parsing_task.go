package models

import "time"

type ParsingStatus string

const (
	ParsingStatusQueued     ParsingStatus = "queued"
	ParsingStatusProcessing ParsingStatus = "processing"
	ParsingStatusCompleted  ParsingStatus = "completed"
	ParsingStatusFailed     ParsingStatus = "failed"
)

type ParsingTask struct {
	ID             string        `bson:"_id" json:"task_id"`
	Status         ParsingStatus `bson:"status" json:"status"`
	SpreadsheetID  string        `bson:"spreadsheet_id" json:"spreadsheet_id"`
	RestaurantName string        `bson:"restaurant_name" json:"restaurant_name"`
	MenuID         string        `bson:"menu_id,omitempty" json:"menu_id,omitempty"`
	ErrorMessage   string        `bson:"error_message,omitempty" json:"error,omitempty"`
	RetryCount     int           `bson:"retry_count" json:"retry_count"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}
