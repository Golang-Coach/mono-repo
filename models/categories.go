package models

import "github.com/jinzhu/gorm"

type Categories struct {
	gorm.Model
	Name string
}
