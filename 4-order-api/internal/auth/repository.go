package auth

import (
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
)

type AuthRepository struct {
	DataBase *db.Db
}

func NewAuthRepository(database *db.Db) *AuthRepository {
	return &AuthRepository{
		DataBase: database,
	}
}

func (repo *AuthRepository) Create(userId *user.UserId) (*user.UserId, error) {
	result := repo.DataBase.DB.Create(userId)
	if result.Error != nil {
		return nil, result.Error
	}

	return userId, nil
}

func (repo *AuthRepository) FindByEmail(
	sessionId string,
) (*user.UserId, error) {
	var user user.UserId
	result := repo.DataBase.DB.First(&user, "phone = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
