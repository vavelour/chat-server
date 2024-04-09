package models

type UserListModel struct {
	Usernames []string `db:"username"`
}
