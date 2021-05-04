package application

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-account/domain/service"
	"github.com/smarest/smarest-common/client/resource"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

type UserService struct {
	userService *service.UserService
}

func NewUserService(userService *service.UserService) *UserService {
	return &UserService{userService: userService}
}

func (s *UserService) Get(c *gin.Context) {
	requestBody := resource.LoginRequestResource{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, exception.GetError(exception.CodeValueInvalid))
		return
	}

	user, err := s.userService.Get(&service.UserLoginInfo{
		UserName: requestBody.UserName,
		Password: requestBody.Password,
		Cookie:   requestBody.Cookie,
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
		tokenString, tErr := s.userService.GenerateUserToken(user)
		if tErr != nil {
			log.Print(tErr.Error())
			c.JSON(http.StatusInternalServerError, exception.GetError(exception.CodeSystemError))
		}
		requestBody.Cookie = tokenString
	}

	var fields = c.Query("fields")
	if fields == "" {
		c.JSON(http.StatusOK, gin.H{"user": user, "cookie": requestBody.Cookie})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user.ToSlide(fields), "cookie": requestBody.Cookie})
	}
}

// Decrepcated
func (s *UserService) GetByCookie(c *gin.Context) {
	cookie := c.Params.ByName("cookie")
	if cookie == "" {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "Cookie is invalid"))
		return
	}

	claims, err := s.userService.CheckUserToken(cookie)

	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	if claims.User == nil {
		c.JSON(http.StatusNotFound, exception.CreateError(exception.CodeNotFound, "UserName not found."))
		return
	}

	c.JSON(http.StatusOK, claims.User)
}
