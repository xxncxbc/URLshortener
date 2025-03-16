package auth

import (
	"URLshortener/configs"
	"URLshortener/internal/user"
	"URLshortener/pkg/db"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	InitialEmail    = "a@a.com"
	initialPassword = "password"
)

func prepare() (*AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := AuthHandler{
		AuthService: NewAuthService(userRepo),
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				AccessSecret:  "secret",
				RefreshSecret: "secret",
			},
		},
	}
	return &handler, mock, nil
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mock, mockErr := prepare()
	if mockErr != nil {
		t.Fatal(mockErr)
		return
	}
	pass, _ := bcrypt.GenerateFromPassword([]byte(initialPassword), bcrypt.MinCost)
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow(initialPassword, string(pass))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	data, _ := json.Marshal(&LoginRequest{
		Email:    InitialEmail,
		Password: initialPassword,
	})
	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(wr, req)
	if wr.Result().StatusCode != http.StatusOK {
		t.Errorf("login: expected status 200, got %d", wr.Result().StatusCode)
	}
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := prepare()
	if err != nil {
		t.Fatal(err)
		return
	}
	mock.ExpectQuery("SELECT").WillReturnError(nil)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectCommit()
	data, _ := json.Marshal(&RegisterRequest{
		Email:    InitialEmail,
		Password: initialPassword,
		Name:     "Name",
	})
	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(wr, req)
	if wr.Result().StatusCode != http.StatusCreated {
		t.Errorf("register: expected status 201, got %d", wr.Result().StatusCode)
	}
}
