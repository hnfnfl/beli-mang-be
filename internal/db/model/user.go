package model

type User struct {
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
	Role         string `db:"role"`
	EmailRole    string `db:"email_role"`
}
