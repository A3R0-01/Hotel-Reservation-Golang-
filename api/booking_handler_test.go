package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelapi.com/db/fixtures"
	"hotelapi.com/types"
)

func TestGetBookings(t *testing.T) {
	var (
		testDB          = setup(t)
		bookingHandler  = NewBookingHandler(testDB.store)
		app             = fiber.New()
		group           = app.Group("/", JWTAuthentication(testDB.store.User))
		testUser        = fixtures.AddUser(testDB.store, false, "james", "foo", types.DefaultUserPassword)
		fakeUser        = fixtures.AddUser(testDB.store, false, "jack", "bauer", types.DefaultUserPassword)
		admin           = fixtures.AddUser(testDB.store, true, "admin", "admin", "1234bsrvnt")
		hotel           = fixtures.AddHotel(testDB.store, "Njeke Hotel", "Harare, Zimbabwe", 4, []primitive.ObjectID{})
		room            = fixtures.AddRoom(testDB.store, 88.50, "large", true, hotel.ID)
		booking         = fixtures.AddBooking(testDB.store, room.ID, testUser.ID, time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 5))
		reqAdmin        = httptest.NewRequest("GET", "/", bytes.NewReader(nil))
		reqUser         = httptest.NewRequest("GET", "/", bytes.NewReader(nil))
		reqFakeUser     = httptest.NewRequest("GET", "/", bytes.NewReader(nil))
		reqNoUser       = httptest.NewRequest("GET", "/", bytes.NewReader(nil))
		token           = CreateTokenFromUser(admin)
		normalUserToken = CreateTokenFromUser(testUser)
		fakeToken       = CreateTokenFromUser(fakeUser)
	)
	defer testDB.teardown(t)
	group.Get("/", bookingHandler.HandleGetBookings)
	reqAdmin.Header.Add("x-api-token", token)
	reqUser.Header.Add("x-api-token", normalUserToken)
	reqFakeUser.Header.Add("x-api-token", fakeToken)
	respAdmin, err := app.Test(reqAdmin)
	if err != nil {
		t.Fatal(err)
	}
	respUser, err := app.Test(reqUser)
	if err != nil {
		t.Fatal(err)
	}
	respFake, err := app.Test(reqFakeUser)
	if err != nil {
		t.Fatal(err)
	}
	respNoUser, err := app.Test(reqNoUser)
	if err != nil {
		t.Fatal(err)
	}
	var bookingsAdmin []*types.Booking
	var bookingsUser []*types.Booking
	var bookingsFake []*types.Booking
	json.NewDecoder(respAdmin.Body).Decode(&bookingsAdmin)
	json.NewDecoder(respUser.Body).Decode(&bookingsUser)
	json.NewDecoder(respFake.Body).Decode(&bookingsFake)

	if respAdmin.Status != "200 OK" {
		t.Fatalf("Request Failed for Admin: %d", respAdmin.StatusCode)
	}
	if respUser.Status != "200 OK" {
		t.Fatalf("Request Failed for User: %d", respUser.StatusCode)
	}
	if respFake.Status != "200 OK" {
		t.Fatalf("Request Failed for Fake User: %d", respUser.StatusCode)
	}
	if respNoUser.StatusCode != http.StatusUnauthorized {
		fmt.Println(respNoUser.StatusCode, http.StatusUnauthorized)
		t.Fatal("expected user to be blocked but instead got access")
	}

	if len(bookingsFake) > 0 {
		t.Fatal("unauthorized access")
	}
	if !(reflect.DeepEqual(booking.ID, bookingsUser[0].ID) && reflect.DeepEqual(bookingsAdmin[0].ID, bookingsUser[0].ID)) {
		t.Fatal("got wrong booking")
	}
	fmt.Println("Pass: Admin:\t", bookingsAdmin, "\t", len(bookingsAdmin))
	fmt.Println("Pass: TestUser:\t", bookingsUser, "\t", len(bookingsUser))
	fmt.Println("Pass: FakeUser:\t", bookingsUser, "\t", len(bookingsFake))
	fmt.Println("Pass: NoUser:\t", respNoUser.Status)
}

func TestGetBooking(t *testing.T) {
	var (
		testDB          = setup(t)
		bookingHandler  = NewBookingHandler(testDB.store)
		app             = fiber.New()
		group           = app.Group("/", JWTAuthentication(testDB.store.User))
		testUser        = fixtures.AddUser(testDB.store, false, "james", "foo", types.DefaultUserPassword)
		fakeUser        = fixtures.AddUser(testDB.store, false, "jack", "bauer", types.DefaultUserPassword)
		admin           = fixtures.AddUser(testDB.store, true, "admin", "admin", "1234bsrvnt")
		hotel           = fixtures.AddHotel(testDB.store, "Njeke Hotel", "Harare, Zimbabwe", 4, []primitive.ObjectID{})
		room            = fixtures.AddRoom(testDB.store, 88.50, "large", true, hotel.ID)
		booking         = fixtures.AddBooking(testDB.store, room.ID, testUser.ID, time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 5))
		hexString       = booking.ID.Hex()
		reqAdmin        = httptest.NewRequest("GET", fmt.Sprint("/", hexString), bytes.NewReader(nil))
		reqUser         = httptest.NewRequest("GET", fmt.Sprint("/", hexString), bytes.NewReader(nil))
		reqFakeUser     = httptest.NewRequest("GET", fmt.Sprint("/", hexString), bytes.NewReader(nil))
		reqNoUser       = httptest.NewRequest("GET", fmt.Sprint("/", hexString), bytes.NewReader(nil))
		token           = CreateTokenFromUser(admin)
		normalUserToken = CreateTokenFromUser(testUser)
		fakeToken       = CreateTokenFromUser(fakeUser)
	)
	defer testDB.teardown(t)
	group.Get("/:id", bookingHandler.HandleGetBooking)
	reqAdmin.Header.Add("x-api-token", token)
	reqUser.Header.Add("x-api-token", normalUserToken)
	reqFakeUser.Header.Add("x-api-token", fakeToken)
	respAdmin, err := app.Test(reqAdmin)
	if err != nil {
		t.Fatal(err)
	}
	respUser, err := app.Test(reqUser)
	if err != nil {
		t.Fatal(err)
	}
	respFake, err := app.Test(reqFakeUser)
	if err != nil {
		t.Fatal(err)
	}
	respNoUser, err := app.Test(reqNoUser)
	if err != nil {
		t.Fatal(err)
	}
	var bookingAdmin types.Booking
	var bookingUser types.Booking
	json.NewDecoder(respAdmin.Body).Decode(&bookingAdmin)
	json.NewDecoder(respUser.Body).Decode(&bookingUser)

	if respAdmin.Status != "200 OK" {
		t.Fatalf("Request Failed for Admin: %d", respAdmin.StatusCode)
	}
	if respUser.Status != "200 OK" {
		t.Fatalf("Request Failed for User: %d", respUser.StatusCode)
	}
	if respFake.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Request Failed for Fake User: %d", respUser.StatusCode)
	}
	if respNoUser.StatusCode != http.StatusUnauthorized {
		fmt.Println(respNoUser.StatusCode, http.StatusUnauthorized)
		t.Fatal("expected user to be blocked but instead got access")
	}

	if !(reflect.DeepEqual(booking.ID, bookingUser.ID) && reflect.DeepEqual(bookingAdmin.ID, bookingUser.ID)) {
		t.Fatal("got wrong booking")
	}
	fmt.Println("Pass: Admin:\t", bookingAdmin.ID)
	fmt.Println("Pass: TestUser:\t", bookingUser.ID)
	fmt.Println("Pass: FakeUser:\t", respFake.Status)
	fmt.Println("Pass: NoUser:\t", respNoUser.Status)
}
