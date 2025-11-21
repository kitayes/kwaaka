package models

type MenuParsingMessage struct {
	TaskID         string `json:"task_id"`
	SpreadsheetID  string `json:"spreadsheet_id"`
	RestaurantName string `json:"restaurant_name"`
}
