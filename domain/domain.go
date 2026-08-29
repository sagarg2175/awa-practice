package domain

import (
	"time"
)

type Master struct {
	Id              *int       `json:"id" db:"id"`
	Name            *string    `json:"name" db:"name"`
	LastUpdatedDate *time.Time `json:"lastupdated" db:"lastupdated"`
	Domain        *string    `json:"domain" db:"domain"`
}
