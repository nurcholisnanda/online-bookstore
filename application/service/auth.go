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
	ValidateToken(string) (uint, error)
}

// Authentication client constructor
func NewAuthClient(secret string, userRepo user.Repository) Authentication {
	return &authClient{
		secret:   secret,
		userRepo: userRepo,
	}
}

// CreateAccessToken will create access token that will be used for user authentication.
// access token will be needed in API that needs user to be authorized
func (c *authClient) CreateAccessToken(userID uint) (accessToken string, err error) {
	if userID == 0 {
		return "", errUserIDRequired
	}
	currentTime := time.Now().UTC()
	//set token with criteria below and input userID into subject
	//this will be needed to check which user is this token for
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

// ValidateTokens will validate whether the token is valid
// and will return user id if the user is exist in our database
// otherwise it will return error
func (c *authClient) ValidateToken(requestToken string) (uint, error) {
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
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errClaimingToken
	}

	if !ok && !token.Valid {
		return 0, errInvalidToken
	}

	currentTime := time.Now().UTC().Unix()
	if ok := claims.VerifyExpiresAt(currentTime, true); !ok {
		return 0, errTokenExpired
	}

	//claim our user id input in subject from token
	userID, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return 0, err
	}

	// Check whether the user is exist
	_, err = c.userRepo.FindUserByID(uint(userID))
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}
