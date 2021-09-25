package main

type JSON struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type User struct {
	ID          uint64   `json:"-"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"-"`
	Age         uint     `json:"age"`
	Description string   `json:"description"`
	ImgSrc      string   `json:"imgSrc"`
	Tags        []string `json:"tags"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
