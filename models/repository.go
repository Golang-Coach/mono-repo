package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Repository struct {
	gorm.Model
	Name          string
	Owner         string
	FullName      string
	Description   string
	Stars         int
	Forks         int
	LastUpdatedBy string
	ReadMe        string       `sql:"type:text;"`
	Tags          []Categories `gorm:"many2many:repo_categories;"`
	Categories    []Tags       `gorm:"many2many:repo_tags;"`
	User          User
	UserID        uint
	Processed     bool
	ProcessedAt   time.Time
}
