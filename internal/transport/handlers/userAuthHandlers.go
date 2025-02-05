package handlers

import (
	"errors"
	"fmt"
	"webPractice1/internal/domain"

	"github.com/gin-gonic/gin"
)

// AuthSignUp godoc
// @Summary Create a new user
// @Description Create a new user and store it in the database
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body domain.User true "Domain User"
// @Success 200 {object} map[string]int "User ID"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /auth/signUp [post]
func (har *HandlerAssetsResponse) singUp(c *gin.Context) {
	var input domain.User
	if err := c.BindJSON(&input); err != nil {
		har.Logger.Error(fmt.Sprintf("BindJSON method error: %s", err))
		c.JSON(400, gin.H{"error": "Invalid data input"})
		return
	}
	id, err := har.service.Autherization.CreateUser(input)
	if err != nil {
		c.JSON(500, gin.H{"error": "something wrong"})
		return
	}
	c.JSON(200, map[string]int{
		"id": id,
	})
}

// AuthSignIn godoc
// @Summary User Authentication
// @Description Authenticates a user using provided credentials (username and password). Returns a JWT token and a refresh token upon successful authentication.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body domain.UserSignIn true "User credentials for sign in"
// @Success 200 {object} map[string]string "JWT token and refresh token"
// @Failure 400 {object} map[string]string "Invalid input or incorrect username/password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/signIn [post]
func (har *HandlerAssetsResponse) singIn(c *gin.Context) {
	var input *domain.UserSignIn
	if err := c.BindJSON(&input); err != nil {
		har.Logger.Error(fmt.Sprintf("BindJSON method error: %s", err))
		c.JSON(400, gin.H{"error": "Invalid data input"})
		return
	}
	token, refreshToken, err := har.service.Session.GenTokens(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(400, gin.H{"error": "user or password is incorrect"})
			return
		}
		c.JSON(500, gin.H{"error": "something wrong"})
		return
	}
	//c.Header("Authorization", "Bearer "+token) //отправить в хедере
	c.SetCookie("CookieToken", refreshToken, 5184000, "/", "localhost", false, true)
	c.JSON(200, map[string]string{
		"token":        token,
		"RefreshToken": refreshToken,
	})
}

// RefreshHandler godoc
// @Summary Refresh Tokens
// @Description Refreshes the JWT token and refresh token using the provided refresh token from the cookie.
// @Tags Auth
// @Accept json
// @Produce json
// @Param CookieToken header string true "Refresh token from the cookie"
// @Success 200 {object} map[string]string "Updated JWT token and new refresh token"
// @Failure 400 {object} map[string]string "Bad cookie or token refresh error"
// @Failure 401 {object} map[string]string "Obsolete token or unauthorized access"
// @Router /auth/refresh [get]
func (har *HandlerAssetsResponse) RefreshHandler(c *gin.Context) {
	oldRefhreshToken, err := c.Cookie("CookieToken")
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Cookie"})
		return
	}
	newJwtToken, newRefreshToken, err := har.service.Session.UpdateTokens(oldRefhreshToken)
	if err != nil {
		if err == domain.ErrObsoleteToken {
			c.JSON(401, gin.H{"error": domain.ErrObsoleteToken.Error()})
			return
		}
		c.JSON(400, gin.H{"error": "Something wrong"})
		return
	}
	c.SetCookie("CookieToken", newRefreshToken, 5184000, "/", "localhost", false, true)
	c.JSON(200, map[string]string{
		"token":        newJwtToken,
		"RefreshToken": newRefreshToken,
	})
}
