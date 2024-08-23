package model

import (
	"github.com/google/uuid"
)

type Instance struct {
	Id       uuid.UUID
	JsonData string
}
