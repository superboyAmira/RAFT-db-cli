package model

import (
	"github.com/google/uuid"
)

type Instance struct {
	Id      uuid.UUID
	Content JsonData
	Term    int64
}

type JsonData struct {
	Name string `json:"name"`
}
