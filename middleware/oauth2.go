package middleware

import (
	"LucasApi/api/config"
	models "LucasApi/api/models"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Register UserInfo for gob encoding (required for sessions)
func init() {
	gob.Register(&models.UserInfo{})
	gob.Register(time.Time{})
}

// OAuth2Handler handles OAuth2 authentication
type OAuth2Handler struct {
	config      *config.OAuth2Config
	oauthConfig *oauth2.Config
}

// NewOAuth2Handler creates a new OAuth2 handler
func NewOAuth2Handler(cfg *config.OAuth2Config) (*OAuth2Handler, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("oauth2 is not enabled")
	}

	var endpoint oauth2.Endpoint

	switch cfg.Provider {
	case "google":
		endpoint = google.Endpoint
		if cfg.AuthURL == "" {
			cfg.AuthURL = google.Endpoint.AuthURL
		}
		if cfg.TokenURL == "" {
			cfg.TokenURL = google.Endpoint.TokenURL
		}
		if cfg.UserInfoURL == "" {
			cfg.UserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
		}
	case "github":
		endpoint = github.Endpoint
		if cfg.AuthURL == "" {
			cfg.AuthURL = github.Endpoint.AuthURL
		}
		if cfg.TokenURL == "" {
			cfg.TokenURL = github.Endpoint.TokenURL
		}
		// Set GitHub user info URL explicitly
		cfg.UserInfoURL = "https://api.github.com/user"
		// GitHub doesn't use "openid" scope, use specific GitHub scopes
		if len(cfg.Scopes) == 0 || containsScope(cfg.Scopes, "openid") {
			cfg.Scopes = []string{"read:user", "user:email"}
		}
	default:
		// Use custom endpoints
		endpoint = oauth2.Endpoint{
			AuthURL:  cfg.AuthURL,
			TokenURL: cfg.TokenURL,
		}
	}

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
		Endpoint:     endpoint,
	}

	return &OAuth2Handler{
		config:      cfg,
		oauthConfig: oauthConfig,
	}, nil
}

// RequireAuth middleware requires authentication
func (h *OAuth2Handler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":     "Unauthorized",
				"message":   "Authentication required",
				"login_url": "/api/auth/login",
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user", user)
		c.Next()
	}
}

// OptionalAuth middleware adds user info to context if available
func (h *OAuth2Handler) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		user := session.Get("user")
		if user != nil {
			c.Set("user", user)
		}

		c.Next()
	}
}

// LoginHandler initiates the OAuth2 login flow
func (h *OAuth2Handler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// Generate state parameter for CSRF protection
		state := generateRandomString(32)
		session.Set("oauth_state", state)
		if err := session.Save(); err != nil {
			log.Printf("Error saving oauth state: %v", err)
		}

		// Get redirect URL parameter
		redirectTo := c.Query("redirect_to")
		if redirectTo != "" {
			session.Set("redirect_after_login", redirectTo)
			if err := session.Save(); err != nil {
				log.Printf("Error saving redirect URL: %v", err)
			}
		}

		url := h.oauthConfig.AuthCodeURL(state)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// CallbackHandler handles the OAuth2 callback
func (h *OAuth2Handler) CallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// Verify state parameter
		state := c.Query("state")
		sessionState := session.Get("oauth_state")
		if sessionState == nil || state != sessionState.(string) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid state parameter",
			})
			return
		}

		// Exchange authorization code for token
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Authorization code not provided",
			})
			return
		}

		token, err := h.oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to exchange code for token",
				"details": err.Error(),
			})
			return
		}

		// Get user information
		userInfo, err := h.getUserInfo(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get user information",
				"details": err.Error(),
			})
			return
		}

		// Store user in session
		session.Set("user", userInfo)
		session.Set("access_token", token.AccessToken)
		// Don't store token expiry in session to avoid gob encoding issues
		session.Delete("oauth_state")
		err = session.Save()
		if err != nil {
			fmt.Printf("DEBUG: Failed to save session: %v\n", err)
		} else {
			fmt.Printf("DEBUG: Session saved successfully with user: %+v\n", userInfo)
		}

		// Check if this is an API request (JSON expected) vs browser request (redirect expected)
		acceptHeader := c.GetHeader("Accept")
		isAPIRequest := strings.Contains(acceptHeader, "application/json") || c.Query("format") == "json"

		redirectTo := session.Get("redirect_after_login")
		if redirectTo != nil {
			session.Delete("redirect_after_login")
			if err := session.Save(); err != nil {
				log.Printf("Error clearing redirect URL: %v", err)
			}

			if isAPIRequest {
				// For API requests, return JSON with redirect URL
				c.JSON(http.StatusOK, gin.H{
					"message":     "Login successful",
					"user":        userInfo,
					"redirect_to": redirectTo.(string),
				})
			} else {
				// For browser requests, redirect as before
				c.Redirect(http.StatusTemporaryRedirect, redirectTo.(string))
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"user":    userInfo,
			})
		}
	}
}

// LogoutHandler handles user logout
func (h *OAuth2Handler) LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		if err := session.Save(); err != nil {
			log.Printf("Error clearing session during logout: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	}
}

// MeHandler returns current user information
func (h *OAuth2Handler) MeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authenticated",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

// getUserInfo fetches user information from the OAuth provider
func (h *OAuth2Handler) getUserInfo(token *oauth2.Token) (*models.UserInfo, error) {
	client := h.oauthConfig.Client(context.Background(), token)

	resp, err := client.Get(h.config.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d", resp.StatusCode)
	}

	var rawUserInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawUserInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	userInfo := &models.UserInfo{
		Provider: h.config.Provider,
	}

	// Map provider-specific fields to our UserInfo struct
	switch h.config.Provider {
	case "google":
		userInfo.ID = getStringFromMap(rawUserInfo, "id")
		userInfo.Name = getStringFromMap(rawUserInfo, "name")
		userInfo.Email = getStringFromMap(rawUserInfo, "email")
		userInfo.Picture = getStringFromMap(rawUserInfo, "picture")
	case "github":
		userInfo.ID = fmt.Sprintf("%v", rawUserInfo["id"])
		userInfo.Name = getStringFromMap(rawUserInfo, "name")
		if userInfo.Name == "" {
			userInfo.Name = getStringFromMap(rawUserInfo, "login") // Fallback to login name
		}
		userInfo.Email = getStringFromMap(rawUserInfo, "email")
		// GitHub may not return email in the user endpoint if it's private
		// Fetch emails separately if not present
		if userInfo.Email == "" {
			if email := h.getGitHubPrimaryEmail(client); email != "" {
				userInfo.Email = email
			}
		}
		userInfo.Picture = getStringFromMap(rawUserInfo, "avatar_url")
	default:
		// Generic mapping for custom providers
		userInfo.ID = getStringFromMap(rawUserInfo, "id")
		userInfo.Name = getStringFromMap(rawUserInfo, "name")
		userInfo.Email = getStringFromMap(rawUserInfo, "email")
		userInfo.Picture = getStringFromMap(rawUserInfo, "picture")
	}

	return userInfo, nil
}

// getGitHubPrimaryEmail fetches the primary email from GitHub's emails API
func (h *OAuth2Handler) getGitHubPrimaryEmail(client *http.Client) string {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		log.Printf("Failed to fetch GitHub emails: %v", err)
		return ""
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch GitHub emails, status: %d", resp.StatusCode)
		return ""
	}

	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		log.Printf("Failed to decode GitHub emails: %v", err)
		return ""
	}

	// Find the primary email
	for _, emailData := range emails {
		if primary, ok := emailData["primary"].(bool); ok && primary {
			if email, ok := emailData["email"].(string); ok {
				return email
			}
		}
	}

	// If no primary email, return the first verified email
	for _, emailData := range emails {
		if verified, ok := emailData["verified"].(bool); ok && verified {
			if email, ok := emailData["email"].(string); ok {
				return email
			}
		}
	}

	return ""
}

// Helper functions
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// containsScope checks if a slice contains a specific scope
func containsScope(scopes []string, scope string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}
