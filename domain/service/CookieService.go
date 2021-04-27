package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/smarest/smarest-account/domain/entity"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

// CookieService get user info
type CookieService struct {
	JWTKey      []byte
	ExpiredTime time.Duration
}

// Claims store userName and jwt keys
type UserClaims struct {
	User *entity.User
	jwt.StandardClaims
}

// NewCookieService is CookieService's constructor
func NewCookieService(key string) *CookieService {
	return &CookieService{
		JWTKey:      []byte(key),
		ExpiredTime: 60 * 24 * time.Minute,
	}
}

// GenerateToken create token
func (s *CookieService) GenerateUserToken(user *entity.User) (string, error) {
	expirationTime := time.Now().Add(s.ExpiredTime)
	claims := &UserClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	return token.SignedString(s.JWTKey)
}

// CheckToken check token
func (s *CookieService) CheckUserToken(cookie string) (*UserClaims, *exception.Error) {
	// Initialize a new instance of `Claims`
	claims := &UserClaims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return s.JWTKey, nil
	})

	if err != nil {
		return nil, exception.CreateError(exception.CodeSignatureInvalid, "Cookie is invalid")
	}
	if !tkn.Valid {
		return nil, exception.CreateError(exception.CodeSignatureInvalid, "Cookie is invalid")
	}
	return claims, nil
}
