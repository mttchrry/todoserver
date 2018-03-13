package models

import (
	"errors"
	"fmt"
	"os"
	"time"
	"todo/src/db"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User - the user creating todos
type User struct {
	ID                string  `db:"id"                  json:"id"`
	Email             string  `db:"email"               json:"email"`
	FirstName         string  `db:"first_name"          json:"first_name"`
	LastName          string  `db:"last_name"           json:"last_name"`
	EncryptedPassword *string `db:"encrypted_password"  json:"-"`
	Password          *string `db:"-"                   json:"password"`
	JwtToken          string  `db:"-"                   json:"jwt_token"`
}

func (u *User) Create() error {
	fmt.Println("User ", *u)
	if u.Password != nil {
		// would need better check of strength for real
		if len(*u.Password) < 6 {
			return errors.New("Password needs to be at least six characters long")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Here2")
			return err
		}

		encryptedPassword := string(hash)
		u.EncryptedPassword = &encryptedPassword
	}

	err := db.Todo.Get(&u.ID, `
	   INSERT INTO users (
		 email,
		 first_name,
		 last_name,
		 encrypted_password
	   ) VALUES ($1, $2, $3, $4) RETURNING id
	 `,
		u.Email,
		u.FirstName,
		u.LastName,
		u.EncryptedPassword,
	)
	fmt.Println("Here3, err ", err)
	GetTokenHandler(u)
	return err
}

// Update syncs the struct instance changes into the database
func (u *User) Update() error {
	var prevUser User
	err := db.Todo.Get(&prevUser,
		`Select * from users where id = $1 or email = $2`,
		u.ID, u.Email)
	if err != nil {
		return errors.New("No user with specified ID to update")
	}

	if u.Password != nil && *u.Password != "" {
		if len(*u.Password) < 6 {
			return errors.New("Password needs to be at least four characters long")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Here2")
			return err
		}

		encryptedPassword := string(hash)
		u.EncryptedPassword = &encryptedPassword
	}

	_, err = db.Todo.Exec(
		`UPDATE users SET
		email=$2,
		first_name=$3,
		last_name=$4,
		encrypted_password=$5
	   WHERE id=$1`,
		prevUser.ID,
		u.Email,
		u.FirstName,
		u.LastName,
		u.EncryptedPassword)
	if err != nil {
		fmt.Printf("Here1")
		return err
	}
	return nil
}

// Update syncs the struct instance changes into the database
func (u *User) Get() error {
	// could potentially get other ways as well, keeping this simple
	if u.Email != "" {
		err := db.Todo.Get(u,
			`Select * from users where email = $1`,
			u.Email)
		if err != nil {
			return errors.New("No user with specified email to get")
		}
	} else if u.ID != "" {
		err := db.Todo.Get(u,
			`Select * from users where id = $1`,
			u.ID)
		if err != nil {
			return errors.New("No user with specified email to get")
		}
	} else {
		return errors.New("Need email or id for user to get")
	}
	return nil
}

// Delete the struct user from the database
func (u *User) Delete() error {
	_, err := db.Todo.Exec(
		`DELETE FROM users where id = $1 or email = $2`,
		u.ID, u.Email)
	if err != nil {
		return errors.New("No user with specified ID or email to delete")
	}
	return nil
}

// Authenticate returns true if the provided password matches the one stored in the database
func (u *User) Authenticate(password string) error {
	if u.EncryptedPassword == nil {
		return errors.New("OAuth user cannot be authenticated with pasword")
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(*u.EncryptedPassword),
		[]byte(password))

	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}

/* Handlers */
var GetTokenHandler = func(user *User) {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["first_name"] = user.FirstName
	claims["last_name"] = user.LastName
	claims["encrypted_password"] = user.EncryptedPassword
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	user.JwtToken = tokenString
}
