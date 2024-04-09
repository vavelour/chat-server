package models

type UserModel struct {
	Username string `db:"username"`
	Password string `db:"password_hash"`
}
