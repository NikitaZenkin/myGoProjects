package user

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
)

func TestUserExit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	repo := NewUserRepo(db)
	user := User{
		ID:       "5",
		UserName: "Nikita",
		password: "nope",
	}
	rows := sqlmock.NewRows([]string{"id", "username", "password"})
	rows.AddRow(user.ID, user.UserName, user.password)
	mock.
		ExpectQuery("SELECT ID, userName, password FROM users WHERE").
		WithArgs("Nikita").
		WillReturnRows(rows)

	item, err := repo.UserExist("Nikita")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(&user, item) {
		t.Errorf("results not match, want %v, have %v", user, item)
		return
	}
	mock.
		ExpectQuery("SELECT ID, userName, password FROM users WHERE").
		WithArgs("Atikin").
		WillReturnError(fmt.Errorf("no rows in result set"))
	_, err = repo.UserExist("Atikin")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	repo := NewUserRepo(db)
	sqlmock.AnyArg()
	user := User{
		ID:       "5",
		UserName: "Nikita",
		password: "nope",
	}
	mock.
		ExpectQuery("INSERT INTO users (`ID`,`userName`, `password`) VALUES").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("db_error"))
	err = repo.AddUser(&user)
	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("%s", err)
	}
}
