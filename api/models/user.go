package models

type UserModel struct {
	ID          string `json:"id" mapstructure:"id"`
	Name        string `json:"name" mapstructure:"name"`
	Email       string `json:"email" mapstructure:"email"`
	Image       string `json:"image" mapstructure:"image"`
	Background  string `json:"background" mapstructure:"background"`
	Description string `json:"description" mapstructure:"description"`
	Birthday    string `json:"birthday" mapstructure:"birthday"`
	Location    string `json:"location" mapstructure:"location"`
}
