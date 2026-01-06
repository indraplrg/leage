package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetCookie(c *gin.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil || cookie == "" {
		logrus.WithError(err)
		return "", errors.New("cookie not found")
	}
	
	return cookie, nil
}