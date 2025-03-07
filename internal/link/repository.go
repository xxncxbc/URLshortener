package link

import (
	"URLshortener/pkg/db"
	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	database *db.Db
}

func NewLinkRepository(db *db.Db) *LinkRepository {
	return &LinkRepository{database: db}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.database.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}
func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.database.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}
