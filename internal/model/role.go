package model

type Role struct {
	ID   uint32 `sql:"primary_key" alias:"id"`
	Name string `alias:"name"`

	Users []UserAccount `alias:"user_account"`
}
