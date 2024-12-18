package models

import (
	"github.com/kacperhemperek/twitter-v2/lib/dbmap"
)

type UserModel struct {
	ID          string           `json:"id" dbmap:"id"`
	Name        string           `json:"name" dbmap:"name"`
	Email       string           `json:"email" dbmap:"email"`
	Image       string           `json:"image" dbmap:"image"`
	Background  *dbmap.NilString `json:"background" dbmap:"background"`
	Description *dbmap.NilString `json:"description" dbmap:"description"`
	Birthday    *dbmap.NilTime   `json:"birthday" dbmap:"birthday"`
	Location    *dbmap.NilString `json:"location" dbmap:"location"`
}
