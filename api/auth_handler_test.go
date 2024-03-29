package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"hotelapi.com/db/fixtures"
	"hotelapi.com/types"
)

func TestAuthenticateSuccess(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)
	app := fiber.New()
	authHandler := NewAuthHandler(testDB.store)
	app.Post("/auth", authHandler.HandleAuthenticate)
	testUser := fixtures.AddUser(testDB.store, false, "james", "foo", types.DefaultUserPassword)

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
	authHandler := NewAuthHandler(testDB.store)
	app.Post("/auth", authHandler.HandleAuthenticate)
	testUser := fixtures.AddUser(testDB.store, false, "james", "foo", types.DefaultUserPassword)
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
