package service

import (
	"log"

	"github.com/smarest/smarest-account/domain/entity"
	"github.com/smarest/smarest-account/domain/repository"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

type UserLoginInfo struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Cookie   string `json:"cookie"`
}

// UserService get user info
type UserService struct {
	*CookieService
	UserRepository repository.UserRepository
}

// NewUserService create UserService
func NewUserService(jwtKey string, userRepository repository.UserRepository) *UserService {
	return &UserService{NewCookieService(jwtKey), userRepository}
}

// GetUser get user information
func (s *UserService) Get(loginInfo *UserLoginInfo) (*entity.User, *exception.Error) {
	if loginInfo.Cookie != "" {
		return s.GetByCookie(loginInfo.Cookie)
	} else if loginInfo.UserName != "" && loginInfo.Password != "" {
		return s.GetByUserNameAndPassword(loginInfo.UserName, loginInfo.Password)
	}
	return nil, exception.GetError(exception.CodeValueInvalid)
}

// GetUserByUserNameAndPassword get user information by userName and password
func (s *UserService) GetByUserNameAndPassword(userName string, password string) (*entity.User, *exception.Error) {
	user, err := s.UserRepository.FindByUserNameAndPassword(userName, password)
	if err != nil {
		log.Print(err.Error())
		return nil, exception.CreateError(exception.CodeNotFound, "UserName or password is invalid.")
	}

	if user == nil {
		log.Print("user is null")
		return nil, exception.CreateError(exception.CodeNotFound, "UserName or password is invalid.")
	}

	return user, nil
}

// GetUserByCookie get user information by cookie
func (s *UserService) GetByCookie(cookie string) (*entity.User, *exception.Error) {
	claims, err := s.CheckUserToken(cookie)
	if err != nil {
		log.Print(err.ErrorMessage)
		return nil, err
	}
	if claims.User == nil {
		resErr := exception.CreateError(exception.CodeValueInvalid, "Cookie is invalid.")
		log.Print(resErr.ErrorMessage)
		return nil, resErr
	}
	return claims.User, nil
}
