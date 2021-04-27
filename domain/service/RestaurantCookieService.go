package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/smarest/smarest-account/domain/entity"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

// CookieService get user info
type RestaurantCookieService struct {
	JWTKey      []byte
	ExpiredTime time.Duration
}

// Claims store restaurantID and jwt keys
type RestaurantClaims struct {
	Restaurant *entity.Restaurant
	jwt.StandardClaims
}

// NewCookieService is CookieService's constructor
func NewRestaurantCookieService(key string) *RestaurantCookieService {
	return &RestaurantCookieService{
		JWTKey:      []byte(key),
		ExpiredTime: 60 * 24 * time.Minute,
	}
}

// GenerateToken create token
func (s *RestaurantCookieService) GenerateRestaurantToken(restaurant *entity.Restaurant) (string, error) {
	expirationTime := time.Now().Add(s.ExpiredTime)
	claims := &RestaurantClaims{
		Restaurant: restaurant,
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
func (s *RestaurantCookieService) CheckRestaurantToken(cookie string) (*RestaurantClaims, *exception.Error) {
	// Initialize a new instance of `Claims`
	claims := &RestaurantClaims{}
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
