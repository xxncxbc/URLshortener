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

func TestLoginSuccess(t *testing.T) {
	const (
		InitialEmail    = "a@a.com"
		initialPassword = "password"
	)
	handler, mock, mockErr := prepare()
	if mockErr != nil {
		t.Fatal(mockErr)
		return
	}
	pass, _ := bcrypt.GenerateFromPassword([]byte(initialPassword), bcrypt.MinCost)
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("a2a.com", string(pass))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	data, _ := json.Marshal(&LoginRequest{
		Email:    "a@a.com",
		Password: "password",
	})
	reader := bytes.NewReader(data)
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(wr, req)
	if wr.Result().StatusCode != http.StatusOK {
		t.Errorf("login: expected status 200, got %d", wr.Result().StatusCode)
	}
}
