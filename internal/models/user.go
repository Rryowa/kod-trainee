package models

//TODO Authorization
//TODO Authentication!!!

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}