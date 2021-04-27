package application

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-account/domain/service"
	"github.com/smarest/smarest-common/client/resource"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

type RestaurantService struct {
	restaurantService *service.RestaurantService
}

func NewRestaurantService(restaurantService *service.RestaurantService) *RestaurantService {
	return &RestaurantService{restaurantService: restaurantService}
}

func (s *RestaurantService) Get(c *gin.Context) {
	requestBody := resource.RestaurantRequestResource{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, exception.GetError(exception.CodeValueInvalid))
		return
	}

	restaurant, err := s.restaurantService.Get(&service.RestaurantLoginInfo{
		ID:        requestBody.ID,
		AccessKey: requestBody.AccessKey,
		Cookie:    requestBody.Cookie,
	})

	if err != nil {
		if err.ErrorCode == exception.CodeSignatureInvalid {
			c.JSON(http.StatusUnauthorized, err)
			return
		}
		if err.ErrorCode == exception.CodeNotFound {
			c.JSON(http.StatusNotFound, err)
			return
		}
		if err.ErrorCode == exception.CodeValueInvalid {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if requestBody.Cookie == "" {
		tokenString, tErr := s.restaurantService.GenerateRestaurantToken(restaurant)
		if tErr != nil {
			log.Print(tErr.Error())
			c.JSON(http.StatusInternalServerError, exception.GetError(exception.CodeSystemError))
		}
		requestBody.Cookie = tokenString
	}

	var fields = c.Query("fields")
	if fields == "" {
		c.JSON(http.StatusOK, gin.H{"restaurant": restaurant, "cookie": requestBody.Cookie})
	} else {
		c.JSON(http.StatusOK, gin.H{"restaurant": restaurant.ToSlide(fields), "cookie": requestBody.Cookie})
	}
}

func (s *RestaurantService) GetByCookie(c *gin.Context) {
	cookie := c.Params.ByName("cookie")
	if cookie == "" {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "Cookie is invalid"))
		return
	}

	claims, err := s.restaurantService.CheckRestaurantToken(cookie)

	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	if claims.Restaurant == nil {
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "UserName not found."))
		return
	}

	c.JSON(http.StatusOK, claims.Restaurant)
}
