package user

import (
	"4-order-api/pkg/db"

	"gorm.io/gorm/clause"
)

type UserRepository struct {
	DataBase *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		DataBase: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.DataBase.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) FindBySession(sessionId string) (*User, error) {
	var user User
	result := repo.DataBase.DB.First(&user, "session_id = ?", sessionId)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *UserRepository) Update(user *User) (*User, error) {
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserRepository) FindByPhone(phone string) (*User, error) {
	var user User
	result := repo.DataBase.DB.First(&user, "phone = ?", phone)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
