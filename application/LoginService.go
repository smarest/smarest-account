package application

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-account/application/resource"
	"github.com/smarest/smarest-account/domain/service"
	"github.com/smarest/smarest-common/domain/entity/exception"
)

const (
	PARAM_FROM_URL = "frm"
)

type LoginService struct {
	HomePage    string
	Domains     []string
	TokenName   string
	UserService *service.UserService
}

func NewLoginService(homePage string, domains []string, tokenName string, userService *service.UserService) *LoginService {
	return &LoginService{HomePage: homePage, Domains: domains, TokenName: tokenName, UserService: userService}
}

func (s *LoginService) GetLogin(c *gin.Context) {
	cookie, err := c.Cookie(s.TokenName)
	fromURL := c.DefaultQuery(PARAM_FROM_URL, s.HomePage)
	resource := &resource.Resource{FromURL: fromURL}
	if err != nil {
		//cookie is NotSet
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"resource": resource})
		return
	}
	_, cErr := s.UserService.GetByCookie(cookie)
	if cErr != nil {
		resource.ErrorMessage = "Cookie is invalid."
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"resource": resource})
		return
	}

	resource.Redirect = fromURL
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"resource": resource})
}

func (s *LoginService) GetLogout(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"resource": &resource.Resource{
			DisableLoader: true,
			FromURL:       c.DefaultQuery(PARAM_FROM_URL, s.HomePage),
			Domains:       s.Domains},
	})
}

func (s *LoginService) PostLogin(c *gin.Context) {
	userName := c.PostForm("user_name")
	password := c.PostForm("password")
	remember := c.PostForm("remember")
	fromURL := c.DefaultPostForm(PARAM_FROM_URL, s.HomePage)
	resource := &resource.Resource{FromURL: fromURL}

	if userName == "" || password == "" {
		//cookie is NotSet
		log.Print(userName, password)
		resource.ErrorMessage = "userName and password is required."
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"resource": resource})
		return
	}
	user, err := s.UserService.GetByUserNameAndPassword(userName, password)
	if err != nil {
		resource.ErrorMessage = err.ErrorMessage
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"resource": resource})
		return
	}

	tokenString, tErr := s.UserService.GenerateUserToken(user)
	if err != nil {
		log.Print(tErr.Error())
		resource.ErrorMessage = exception.GetErrorMessage(exception.CodeSystemError)
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"resource": resource})
	}

	if remember != "" {
		resource.AccessToken = tokenString
		resource.Domains = s.Domains
	}
	resource.Redirect = fromURL
	c.HTML(http.StatusOK, "login.tmpl", gin.H{"resource": resource})
}
