package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteCookie(c *gin.Context, name string) {
	c.SetCookieData(&http.Cookie{
		Name : "paseto_token", 
		Value: "", 
		MaxAge: -1, 
		Path: "/", 
		Domain: "localhost", 
		Secure: false, 
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}