package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelapi.com/db"
	"hotelapi.com/types"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TestDBNAME),
	}
}

func TestPostUser(t *testing.T) {
	testDB := setup(t)
	defer testDB.teardown(t)
	app := fiber.New()
	userHandler := NewUserHandler(testDB.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		Email:     "some@gmail.com",
		FirstName: "james",
		LastName:  "fooo",
		Password:  "hellomanu",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
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
	app := fiber.New()
	userHandler := NewUserHandler(testDB.UserStore)
	app.Get("/", userHandler.HandleGetUsers)
	req := httptest.NewRequest("GET", "/", bytes.NewReader(nil))
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var users []types.User
	json.NewDecoder(resp.Body).Decode(&users)
	if resp.Status != "200 OK" {
		t.Error("Request Failed")
	}
	fmt.Println(users, "\n", resp.Status)

}
