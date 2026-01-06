package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetCookie(c *gin.Context, name, value string, expire int) {
	c.SetCookieData(&http.Cookie{
		Name : name, 
		Value: value, 
		MaxAge: expire, 
		Path: "/", 
		Domain: viper.GetString("host"), 
		Secure: false, 
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}