package main

import (
	"URLshortener/internal/auth"
	"URLshortener/internal/user"
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Create(&user.User{
		Password: string(password),
		Email:    "login@example.com",
		Name:     "John Doe",
	})
}

func removeData(db *gorm.DB) {
	db.
		Unscoped().
		Where("email = ?", "login@example.com").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "login@example.com",
		Password: "password",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected status code %d, got %d", 200, res.StatusCode)
	}
	var loginResponse auth.LoginResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		t.Fatal(err)
	}
	if loginResponse.AccessToken == "" {
		t.Fatalf("Expected access token, got empty")
	}
	if loginResponse.RefreshToken == "" {
		t.Fatalf("Expected refresh token, got empty")
	}
	removeData(db)
}

func TestLoginFail(t *testing.T) {
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "login@example.com",
		Password: "wrong password",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 401 {
		t.Fatalf("Expected status code %d, got %d", 401, res.StatusCode)
	}
	removeData(db)
}
