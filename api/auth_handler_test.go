package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

func TestAuthenticateSuccess(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)
	app := fiber.New()
	authHandler := NewAuthHandler(testDB.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)
	testUser := makeTestUser(t, testDB.UserStore)
	params := AuthParams{
		Email:    testUser.Email,
		Password: types.DefaultUserPassword,
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("response failed")
	}
	var response AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.Token == "" {
		t.Fatal("expected token to be present in the response")
	}
	// set password to empty string because it is not returned in json response
	testUser.EncryptedPassword = ""
	if !reflect.DeepEqual(response.User, testUser) {

		t.Fatal("expected user to be the test user but got a different user")
	}
}
func TestAuthenticateWrongPassword(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)
	app := fiber.New()
	authHandler := NewAuthHandler(testDB.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)
	testUser := makeTestUser(t, testDB.UserStore)
	params := AuthParams{
		Email:    testUser.Email,
		Password: "hellopeople",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatal("response failed")
	}
	var response genericResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s ", response.Type)
	}
}
func makeTestUser(t *testing.T, userStore db.UserStore) *types.User {
	newUser := types.CreateUserParams{
		Email:     "some@gmail.com",
		FirstName: "james",
		LastName:  "fooo",
		Password:  types.DefaultUserPassword,
	}
	encryptedTestUser, err := types.NewUserFromParams(newUser)
	if err != nil {
		t.Fatal(err)

	}
	testUser, err := userStore.CreateUser(context.Background(), encryptedTestUser)
	if err != nil {
		t.Fatal(err)

	}
	return testUser
}
