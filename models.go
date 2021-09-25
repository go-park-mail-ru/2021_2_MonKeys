package main

type User struct {
	ID          uint64
	Name        string
	Email       string
	Password    string
	Age         uint
	Description string
	ImgSrc      string
	Tags        []string
}

var (
	users   = make(map[uint64]User)
	cookies = make(map[string]uint64)
)
