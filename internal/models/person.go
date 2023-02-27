package models

import (
	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm"
)

type Person struct {
	// need to set in db : CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name  string    `gorm:"type:varchar(100)"`
	EMail string    `gorm:"type:varchar(100)"`
	Books []Book
}
