package models

type ParseRequest struct {
	SpreadsheetID  string `json:"spreadsheet_id"`
	RestaurantName string `json:"restaurant_name"`
}
