package dtos

type UserDto struct{
	ID		int		`json:"id"`
	FirstName	string		`json:"first_name`
	LastName	string		`json:"last_name"`
	Email		string		`json:"Email"`
	Password	string		`json:"password"`
	Type		bool		`json:"type"` //True para admin
}

type UsersDto []UserDto