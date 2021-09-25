package main

type MockDB struct {
	//DB int
}

func (MockDB) getUserModel(string) (User, error) {
	return User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "mumeu222@mail.ru",
		Password:    "VBif222!",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/static/users/user1",
		Tags:        []string{"haha", "hihi"},
	}, nil
}


type MockSessionDB struct {

}

func (MockSessionDB) getUserByCookie(string) (User, error) {
	return User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "mumeu222@mail.ru",
		Password:    "VBif222!",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/static/users/user1",
		Tags:        []string{"haha", "hihi"},
	}, nil
}

func (MockSessionDB) newSessionCookie(uint64, string) error {
	return nil
}
