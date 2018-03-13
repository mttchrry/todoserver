package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"todo/src/models"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	password := user.Password

	err = user.Get()
	if err != nil {
		http.Error(w, "EMAIL_NOT_REGISTERED", 500)
		return
	}

	if password == nil || len(*password) == 0 {
		http.Error(w, "PASSWORD_REQUIRED", 500)
		return
	}
	err = user.Authenticate(*password)
	if err != nil {
		http.Error(w, "INCORRECT_PASSWORD", 500)
		return
	}

	models.GetTokenHandler(&user)

	json.NewEncoder(w).Encode(user)
}

/* Set up a global string for our secret */
var mySigningKey = []byte(os.Getenv("JWT_SECRET"))

// JwtMiddleware handler for ensuring access token there for protected functions
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func GetUserFromToken(r *http.Request) (models.User, error) {
	user := models.User{}
	var tokenString = r.Header["Authorization"]
	if len(tokenString) == 0 {
		return user, errors.New("No authorized user")
	}
	fmt.Println(tokenString[0])
	ts := strings.Split(tokenString[0], " ")
	if len(ts) < 2 {
		return user, errors.New("Expecting 'Bearer <tokenstring>' authorization")
	}
	token, err := jwt.Parse(ts[1], func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Don't forget to validate the alg is what you expect:
		if method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("Expected signing method: %v, got %v",
				jwt.SigningMethodHS256.Alg(), method.Alg())
		}

		return mySigningKey, nil
	})
	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.ID = claims["id"].(string)
		fmt.Println("id and name from token", claims["id"], claims["first_name"])
		err := user.Get()
		if err != nil {
			fmt.Println("Couldn't get user from token")
		}
		return user, err
	}
	return user, errors.New("Couldn't validate token")
}
