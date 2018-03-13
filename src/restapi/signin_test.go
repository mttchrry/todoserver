package restapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/src/models"
)

func TestSignIn(t *testing.T) {
	pswd := "test123"
	user := models.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "tester@testing.com",
		Password:  &pswd,
	}
	err := user.Create()
	if err != nil {
		t.Fatal(err)
	}
	defer user.Delete()

	if user.JwtToken == "" {
		t.Fatal("Should be making access token on user create")
	}

	s, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	b := bytes.NewBuffer(s)
	req, err := http.NewRequest("POST", "/signin", b)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SignIn)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	user2 := models.User{}
	err = json.NewDecoder(rr.Body).Decode(&user2)
	if err != nil {
		t.Fatal(err)
	}
	if user2.JwtToken != user.JwtToken {
		t.Errorf("expecting tokens to be equivalent, %s, %s", user.JwtToken, user2.JwtToken)
	}
}
