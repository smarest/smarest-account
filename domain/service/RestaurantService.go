package service

import (
	"log"
	"strconv"

	"github.com/smarest/smarest-account/domain/entity"
	"github.com/smarest/smarest-account/domain/repository"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

type RestaurantLoginInfo struct {
	ID        string `json:"id"`
	AccessKey string `json:"accessKey"`
	Cookie    string `json:"cookie"`
}

type RestaurantService struct {
	*RestaurantCookieService
	RestaurantRepository repository.RestaurantRepository
}

func NewRestaurantService(jwtKey string, restaurantRepository repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{NewRestaurantCookieService(jwtKey), restaurantRepository}
}

func (s *RestaurantService) Get(loginInfo *RestaurantLoginInfo) (*entity.Restaurant, *exception.Error) {
	if loginInfo.Cookie != "" {
		return s.GetByCookie(loginInfo.Cookie)
	} else if loginInfo.ID != "" && loginInfo.AccessKey != "" {
		restIDInt, err := strconv.ParseInt(loginInfo.ID, 10, 64)
		if err != nil {
			log.Print(err.Error())
			return nil, exception.GetError(exception.CodeValueInvalid)
		}
		return s.GetByIDAndAccessKey(restIDInt, loginInfo.AccessKey)
	}
	return nil, exception.GetError(exception.CodeValueInvalid)
}

func (s *RestaurantService) GetByIDAndAccessKey(restaurantID int64, accessKey string) (*entity.Restaurant, *exception.Error) {
	restaurant, err := s.RestaurantRepository.FindByIDAndAccessKey(restaurantID, accessKey)
	if err != nil {
		log.Print(err.Error())
		return nil, exception.CreateError(exception.CodeNotFound, "RestaurantID or AccessKey is invalid.")
	}

	if restaurant == nil {
		log.Print("user is null")
		return nil, exception.CreateError(exception.CodeNotFound, "RestaurantID or AccessKey is invalid.")
	}

	return restaurant, nil
}

func (s *RestaurantService) GetByCookie(cookie string) (*entity.Restaurant, *exception.Error) {
	claims, err := s.CheckRestaurantToken(cookie)
	if err != nil {
		log.Print(err.ErrorMessage)
		return nil, err
	}
	if claims.Restaurant == nil {
		resErr := exception.CreateError(exception.CodeValueInvalid, "Cookie is invalid.")
		log.Print(resErr.ErrorMessage)
		return nil, resErr
	}
	return claims.Restaurant, nil
}
