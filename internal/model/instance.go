package model

import (
	"github.com/google/uuid"
)

type Instance struct {
	Id       uuid.UUID
	JsonData string // here must be a special json struct maybe 
	Err string
}
