package user

import (
	"database/sql"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{DB: db}
}

func (repa *UserRepo) UserExist(userName string) (*User, error) {
	user := &User{}
	err := repa.DB.
		QueryRow("SELECT ID, userName, password FROM users WHERE userName = ?", userName).
		Scan(&user.ID, &user.UserName, &user.password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repa *UserRepo) AddUser(newUser *User) error {
	_, err := repa.DB.Exec(
		"INSERT INTO users (`ID`,`userName`, `password`) VALUES (?, ?, ?)",
		newUser.ID,
		newUser.UserName,
		newUser.password,
	)
	return err
}
