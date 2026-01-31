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

// @Summary Check authentication status
// @Description Returns current authentication status and user info if authenticated
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{} "Authentication status"
// @Router /api/auth/status [get]
func (ac *AuthController) GetAuthStatus(c *gin.Context) {
	// Check if user info is available in context (set by middleware)
	user, exists := c.Get("user")
	if exists && user != nil {
		// Ensure we return all user information including email
		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"user":          user,
		})
		return
	}

	// Check OAuth2 configuration
	oauth2Enabled := c.GetBool("oauth2_enabled")
	oauth2Provider := c.GetString("oauth2_provider")

	c.JSON(http.StatusOK, gin.H{
		"authenticated":  false,
		"message":        "Not authenticated",
		"login_url":      "/api/auth/login",
		"oauth2_enabled": oauth2Enabled,
		"provider":       oauth2Provider,
	})
}
