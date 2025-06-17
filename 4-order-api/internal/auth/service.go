package auth

import (
	"4-order-api/internal/user"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewUserService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{userRepository}
}

func generateSecureSessionID() (string, error) {
	b := make([]byte, 24) // 24 байта дадут 32 символа в base64
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Убираем символы /=+ для чистого SessionID
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

func (service *AuthService) CreateSession(
	phone string,
) (*user.User, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)

	sessionId, err := generateSecureSessionID()
	for {
		if err != nil {
			return nil, err
		}
		userExist, _ := service.UserRepository.FindBySession(sessionId)
		if userExist == nil {
			break
		}
		sessionId, err = generateSecureSessionID()
	}

	if existedUser != nil {
		existedUser.SessionId = sessionId
		_, err := service.UserRepository.Update(existedUser)
		if err != nil {
			return nil, err
		}

		return existedUser, nil
	} else {

		user := &user.User{
			Phone:     phone,
			Code:      3452,
			SessionId: sessionId,
		}

		_, err := service.UserRepository.Create(user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func (service *AuthService) VerifyCode(
	sessionId string, code uint,
) (string, error) {
	existedUser, _ := service.UserRepository.FindBySession(sessionId)
	if existedUser == nil || existedUser.Code != code {
		return "", errors.New(ErrWrongCode)
	}
	return existedUser.Phone, nil
}
