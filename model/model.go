package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
)

type Name struct {
	gorm.Model
	Name string `json:"name"`
}
