package model

type UserAccount struct {
	ID          uint32  `sql:"primary_key" alias:"id"`
	FirstName   string  `alias:"first_name"`
	LastName    string  `alias:"last_name"`
	Email       string  `alias:"email"`
	PhoneNumber *string `alias:"phone_number"`
	Password    string  `alias:"password"`

	Roles []Role `alias:"role"`
}
