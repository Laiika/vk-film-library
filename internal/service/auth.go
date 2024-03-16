package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
	"vk-film-library/internal/entity"
	"vk-film-library/internal/repo"
	"vk-film-library/internal/repo/repoerrs"
)

const (
	salt = "15dd01c7259448d497ec85b125f11bde"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserRole string
}

type AuthService struct {
	userRepo repo.UserRepo
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(userRepo repo.UserRepo, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input *entity.CreateInput) (int, error) {
	user := &entity.User{
		Username: input.Username,
		Password: hash(input.Password),
		Role:     input.Role,
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(ctx context.Context, input *entity.AuthInput) (string, error) {
	// get user from DB
	user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, hash(input.Password))
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserRole: user.Role,
	})

	// sign token
	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return "", ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return "", ErrCannotParseToken
	}

	return claims.UserRole, nil
}

func hash(password string) string {
	h := sha1.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum([]byte(salt)))
}
