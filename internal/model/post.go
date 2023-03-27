package model

type Post struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Text string `json:"text" db:"text"`
}
