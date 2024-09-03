package model

type Instance struct {
	Id      string
	Content JsonData
	Term    int64
}

type JsonData struct {
	Name string `json:"name"`
}
