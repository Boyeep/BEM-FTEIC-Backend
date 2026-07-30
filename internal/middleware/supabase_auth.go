package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

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
			response.Fail(c, apperr.New("UNAUTHORIZED", "missing bearer token", http.StatusUnauthorized))
			c.Abort()
			return
		}

		req, err := http.NewRequestWithContext(
			c.Request.Context(),
			http.MethodGet,
			a.url+"/auth/v1/user",
			nil,
		)
		if err != nil {
			response.Fail(c, apperr.Internal(err))
			c.Abort()
			return
		}
		req.Header.Set("apikey", a.apiKey)
		req.Header.Set("Authorization", "Bearer "+token)

		authResponse, err := a.client.Do(req)
		if err != nil {
			response.Fail(c, apperr.New("AUTH_UNAVAILABLE", "authentication service unavailable", http.StatusBadGateway))
			c.Abort()
			return
		}
		defer authResponse.Body.Close()

		if authResponse.StatusCode != http.StatusOK {
			response.Fail(c, apperr.New("UNAUTHORIZED", "invalid or expired token", http.StatusUnauthorized))
			c.Abort()
			return
		}

		var user supabaseUser
		if err := json.NewDecoder(authResponse.Body).Decode(&user); err != nil || user.ID == "" {
			response.Fail(c, apperr.New("UNAUTHORIZED", "invalid user response", http.StatusUnauthorized))
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Next()
	}
}
