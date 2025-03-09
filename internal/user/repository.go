package user

import (
	"URLshortener/pkg/db"
)

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.database.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	result := repo.database.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
