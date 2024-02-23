package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"hotelapi.com/db/fixtures"
	"hotelapi.com/types"
)

func TestPostUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)
	app := fiber.New()
	userHandler := NewUserHandler(testDB.store)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		Email:     "some@gmail.com",
		FirstName: "james",
		LastName:  "fooo",
		Password:  types.DefaultUserPassword,
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	fmt.Println(user)
	if len(user.ID) == 0 {
		t.Error("Expected a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("Expected not to include a password")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("Expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("Expected email %s but got %s", params.Email, user.Email)
	}
}
func TestGetUsers(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)
	var (
		app         = fiber.New()
		userHandler = NewUserHandler(testDB.store)
		testUser    = fixtures.AddUser(testDB.store, true, "james", "foo", types.DefaultUserPassword)
		req         = httptest.NewRequest("GET", "/", bytes.NewReader(nil))
		group       = app.Group("/", JWTAuthentication(testDB.store.User))

		token = CreateTokenFromUser(testUser)
	)
	group.Get("/", userHandler.HandleGetUsers)
	req.Header.Add("x-api-token", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var users []*types.User
	json.NewDecoder(resp.Body).Decode(&users)
	if resp.Status != "200 OK" {
		t.Fatal("Request Failed")
	}
	fmt.Println(users, "\t", resp.Status)
}
