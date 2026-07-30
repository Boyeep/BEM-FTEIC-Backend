package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SupabaseAuth struct {
	url    string
	apiKey string
	client *http.Client
}

type supabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func NewSupabaseAuth(url, apiKey string) *SupabaseAuth {
	return &SupabaseAuth{
		url:    strings.TrimRight(url, "/"),
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *SupabaseAuth) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(strings.TrimPrefix(
			c.GetHeader("Authorization"),
			"Bearer ",
		))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		req, err := http.NewRequestWithContext(
			c.Request.Context(),
			http.MethodGet,
			a.url+"/auth/v1/user",
			nil,
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "create auth request"})
			return
		}
		req.Header.Set("apikey", a.apiKey)
		req.Header.Set("Authorization", "Bearer "+token)

		response, err := a.client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "authentication service unavailable"})
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		var user supabaseUser
		if err := json.NewDecoder(response.Body).Decode(&user); err != nil || user.ID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user response"})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Next()
	}
}
