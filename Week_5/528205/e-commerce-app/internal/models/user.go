package models

type User struct {
	ID        int    `db:"id"`
	Email     string `db:"email"`
	Password  string `db:"password"` // In a real application, you should hash and salt passwords.
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
