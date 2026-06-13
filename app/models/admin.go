package models

type AdminStats struct {
	Users           int `json:"users"`
	Vehicles        int `json:"vehicles"`
	Recommendations int `json:"recommendations"`
	Questions       int `json:"questions"`
	ActiveUsersWeek int `json:"active_users_week"`
	NewUsersWeek    int `json:"new_users_week"`
}
