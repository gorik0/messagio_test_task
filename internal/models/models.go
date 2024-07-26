package models

type KafkaMesg struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}
type JsonMsg struct {
	Data any `json:"data"`
}

type GetStat struct {
	Total     int `json:"total"`
	Processed int `json:"processed"`
}
