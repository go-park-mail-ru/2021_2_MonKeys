package repository_test

// func getDataBase() (*sqlx.DB, error) {
// 	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

// 	addGetByEmailSupport(mock)
// 	addGetByIdSupport(mock)
// 	addAddSupport(mock)
// 	addUpdateSupport(mock)

// 	sqlxDB := sqlx.NewDb(db, "sqlmock")
// 	return sqlxDB, err
// }

// func TestAdd(t *testing.T) {
// 	lu := models.LoginUser{
// 		Email:    "valid@valid.ru",
// 		Password: "123",
// 	}
// 	con, err := getDataBase()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	rep := repository.NewPostgresStaffRepository(con)
// 	_, err = rep.AddLogin(context.TODO(), lu)
// 	assert.NotNil(t, err)

// }
