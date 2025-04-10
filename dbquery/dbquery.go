package dbquery

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
)

var (
	db          *sql.DB
	NotFoundErr = errors.New("error: not found")
)

func ConnectToDb(user string, pass string, addr string, dbName string) error {
	cfg := mysql.Config{
		User:                 user,
		Passwd:               pass,
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               dbName,
		AllowNativePasswords: true,
		AllowOldPasswords:    true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return err
	}

	err = db.Ping()

	if err != nil {
		return err
	}

	return nil
}

type (
	User struct {
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		PhoneNum  string `json:"phoneNum"`
		Password  string `json:"password"`
		BirthDate string `json:"birthDate"`
	}

	Product struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Amount      int     `json:"amount"`
		Description string  `json:"description"`
	}
)

func RegisterUser(user User) (sql.Result, error) {
	return db.Exec("INSERT INTO users (username, first_name, last_name, email, phone_num, password, birth_date) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Username, user.FirstName, user.LastName, user.Email, user.PhoneNum, user.Password, user.BirthDate)
}

func SelectUserByName(username string) (User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE username LIKE ?", username)

	if err != nil {
		return User{}, err
	}

	var user User

	if !rows.Next() {
		return User{}, NotFoundErr
	}

	err = rows.Scan(&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNum,
		&user.Password,
		&user.BirthDate)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func UpdateUser(newUser User) error {
	username := newUser.Username
	oldUser, err := SelectUserByName(username)

	if err != nil {
		return err
	}

	if newUser.FirstName != oldUser.LastName {
		if _, err = updateUserFirstName(username, newUser.FirstName); err != nil {
			return err
		}
	}

	if newUser.LastName != oldUser.LastName {
		if _, err = updateUserLastName(username, newUser.LastName); err != nil {
			return err
		}
	}

	if newUser.Email != oldUser.Email {
		if _, err = updateUserEmail(username, newUser.Email); err != nil {
			return err
		}
	}

	if newUser.PhoneNum != oldUser.PhoneNum {
		if _, err = updateUserPhoneNum(username, newUser.PhoneNum); err != nil {
			return err
		}
	}

	return nil
}

func updateUserFirstName(username string, firstName string) (sql.Result, error) {
	return db.Exec("UPDATE users SET first_name = ? WHERE username LIKE ?", firstName, username)
}

func updateUserLastName(username string, lastName string) (sql.Result, error) {
	return db.Exec("UPDATE users SET last_name = ? WHERE username LIKE ?", lastName, username)
}

func updateUserEmail(username string, email string) (sql.Result, error) {
	return db.Exec("UPDATE users SET email = ? WHERE username LIKE ?", email, username)
}

func updateUserPhoneNum(username string, phoneNum string) (sql.Result, error) {
	return db.Exec("UPDATE users SET phone_num = ? WHERE username LIKE ?", phoneNum, username)
}

func UpdateUserPassword(username string, password string) (sql.Result, error) {
	return db.Exec("UPDATE users SET password = ? WHERE username LIKE ?", password, username)
}

func SelectAllProducts() ([]Product, error) {
	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	var products []Product

	for rows.Next() {
		var product Product

		err = rows.Scan(&product.Name,
			&product.Price,
			&product.Amount,
			&product.Description)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
