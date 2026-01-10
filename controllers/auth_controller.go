package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related endpoints
type AuthController struct{}

// NewAuthController creates a new auth controller
func NewAuthController() *AuthController {
	return &AuthController{}
}

// @Summary Get login URL
// @Description Returns the OAuth2 login URL with optional redirect parameter
// @Tags auth
// @Produce json
// @Param redirect_to query string false "URL to redirect after login"
// @Success 200 {object} map[string]interface{} "Login information"
// @Router /api/auth/login-url [get]
func (ac *AuthController) GetLoginURL(c *gin.Context) {
	redirectTo := c.Query("redirect_to")

	loginURL := "/api/auth/login"
	if redirectTo != "" {
		loginURL += "?redirect_to=" + redirectTo
	}

	c.JSON(http.StatusOK, gin.H{
		"login_url":   loginURL,
		"message":     "Navigate to this URL to authenticate",
		"redirect_to": redirectTo,
	})
}
