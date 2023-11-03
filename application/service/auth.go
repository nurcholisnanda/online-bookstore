package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nurcholisnanda/online-bookstore/domain/user"
)

var (
	errClaimingToken  = errors.New("claim token error")
	errInvalidToken   = errors.New("invalid token error")
	errTokenExpired   = errors.New("token expired error")
	errUserIDRequired = errors.New("user id is required")
)

type authClient struct {
	secret   string
	userRepo user.Repository
}

//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock
type Authentication interface {
	CreateAccessToken(uint) (string, error)
	ValidateToken(string) (int, error)
}

// NewClient returns a wrapper around authentication client.
func NewAuthClient(secret string, userRepo user.Repository) Authentication {
	return &authClient{
		secret:   secret,
		userRepo: userRepo,
	}
}

func (c *authClient) CreateAccessToken(userID uint) (accessToken string, err error) {
	if userID == 0 {
		return "", errUserIDRequired
	}
	currentTime := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Hour * 12)),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(c.secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func (c *authClient) ValidateToken(requestToken string) (int, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secret), nil
	})

	if err != nil {
		return 0, err
	}

	// assert jwt.MapClaims type
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, errClaimingToken
	}

	if !ok && !token.Valid {
		return 0, errInvalidToken
	}

	currentTime := time.Now().UTC()
	if ok := claims.VerifyExpiresAt(currentTime, true); !ok {
		return 0, errTokenExpired
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	// Check whether the user is exist
	_, err = c.userRepo.FindUserByID(uint(userID))
	if err != nil {
		return 0, err
	}

	return userID, nil
}
