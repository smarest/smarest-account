package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smarest/smarest-account/application"
)

func main() {
	//connect to DB
	bean, err := application.InitBean()
	defer bean.DestroyBean()
	if err != nil {
		log.Fatalln("can not create bean", err)
	}
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/login", bean.LoginService.GetLogin)
	router.GET("/logout", bean.LoginService.GetLogout)
	router.POST("/login", bean.LoginService.PostLogin)

	v1 := router.Group("v1")
	{
		v1.GET("/login/user/:cookie", bean.UserService.GetByCookie)
		v1.POST("/login/user", bean.UserService.Get)
		v1.GET("/login/restaurant/:cookie", bean.RestaurantService.GetByCookie)
		v1.POST("/login/restaurant", bean.RestaurantService.Get)
	}
	router.Run(":8080")
}
