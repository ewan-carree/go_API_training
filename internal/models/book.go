package models

import (
	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
)

type Book struct {
	// need to set in db : CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title    string    `gorm:"type:varchar(100)"`
	Author   string    `gorm:"type:varchar(100)"`
	PersonID uuid.UUID `gorm:"type:uuid;"`
}
