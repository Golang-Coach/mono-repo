package models

import "github.com/jinzhu/gorm"

type Tags struct {
	gorm.Model
	Name string
}
