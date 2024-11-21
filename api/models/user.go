package models

import (
	"github.com/kacperhemperek/twitter-v2/lib/dbmap"
)

type UserModel struct {
	ID          string           `json:"id" dbmap:"id" mapstructure:"id"`
	Name        string           `json:"name" dbmap:"name" mapstructure:"name"`
	Email       string           `json:"email" dbmap:"email" mapstructure:"email"`
	Image       string           `json:"image" dbmap:"image" mapstructure:"image"`
	Background  *dbmap.NilString `json:"background" dbmap:"background" mapstructure:"background"`
	Description *dbmap.NilString `json:"description" dbmap:"description" mapstructure:"description"`
	Birthday    *dbmap.NilTime   `json:"birthday" dbmap:"birthday" mapstructure:"birthday"`
	Location    *dbmap.NilString `json:"location" dbmap:"location" mapstructure:"location"`
}
