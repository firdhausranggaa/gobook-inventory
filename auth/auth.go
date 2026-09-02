package auth

import (
	"book-inventory/models"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func HomeHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func LoginGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"content": "",
	})
}

func LoginPostHandler(c *gin.Context) {
	var credential models.Credentials

	if err := c.Bind(&credential); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Binding error"})
		return
	}

	if credential.Username != os.Getenv("SUPER_USER") || credential.Password != os.Getenv("SUPER_PASS") {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Username/password is invalid"})
	} else {
		claim := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "inventory-book",
			IssuedAt:  time.Now().Unix(),
		}

		sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		token, err := sign.SignedString([]byte(os.Getenv("SUPER_SECRET")))
		
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Token signing error"})
			c.Abort()
			return
		}

		q := url.Values{}
		q.Set("auth", token)
		location := url.URL{Path: "/books", RawQuery: q.Encode()}
		c.Redirect(http.StatusMovedPermanently, location.RequestURI())
	}
}