package respositories

import (
	"database/sql"
	"ryan/app"
	"ryan/util/crypto"
)

//GetUserByEmail gets s user by email
func GetUserByEmail(db *sql.DB, email string) (*app.User, error) {
	const sqlstr = `
	select 
		first_name,
		last_name,
		email,
		password,
		salt
	from users
	where email = :p1
	`
	var user app.User
	err := db.QueryRow(sqlstr, email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Salt)
	return &user, err
}

//AddUserToDatabase creates a new user
func AddUserToDatabase(db *sql.DB, firstName, lastName, email, password string) error {
	const sqlstr = `
	insert into users (
		first_name,
		last_name,
		email,
		password,
		salt
	) values (
		:p1,
		:p2,
		:p3,
		:p4,
		:p5
	)
	`
	salt := crypto.GenerateSalt()
	hashedPassword := crypto.HashPassword(password, salt)
	_, err := db.Exec(sqlstr, firstName, lastName, email, hashedPassword, salt)
	return err
}

//CheckEmailExists is a function
func CheckEmailExists(db *sql.DB, email string) (bool, error) {
	const sqlstr = "select email from users where email=:p1"
	var exists string
	err := db.QueryRow(sqlstr, email).Scan(&exists)
	if exists == "" {
		return false, nil
	}
	return true, err
}
