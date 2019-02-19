package respositories

import (
	"database/sql"

	"ryan/app"
)

//GETUSERBYEMAIL gets s user by email
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
	err := db.QueryRow(sqlstr, email).Scan(&user.First_name, &user.Last_name, &user.Email, &user.Password, &user.Salt)
	return &user, err
}

//AddUserTodatabase creates a new user
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
	err := db.QueryRow(sqlstr, firstName, lastName, email, hashedPassword, salt)
	return err
}

func CheckEmailExists(db *sql.DB, email string) (bool, error) {
	const sqlstr = "select exists(select 1 from users where email= :p1)"
	var exists bool
	err := db.QueryRow(sqlstr, email).Scan(&exists)
	return exists, err
}
