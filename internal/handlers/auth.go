package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/sport-booking/internal/auth"
)

// RegisterAuthRoutes 注册认证路由
func RegisterAuthRoutes(r *gin.Engine, client *auth.Client) {
	if client == nil {
		return
	}

	r.POST("/auth/signup", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		user, err := client.SignUp(body.Email, body.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	r.POST("/auth/login", func(c *gin.Context) {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		session, err := client.SignIn(body.Email, body.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, session)
	})
}
