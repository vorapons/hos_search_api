package ginmiddleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Login      string `json:"login"`
	HospitalID string `json:"hospital_id"`
	Hospital   string `json:"hospital"`
	jwt.RegisteredClaims
}

// JWTProtected validates the Bearer token and injects claims into c context.
// isBlacklisted is called to reject logged-out tokens.
func JWTProtected(jwtSecret string, isBlacklisted func(string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Missing or invalid Authorization header",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if isBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Token has been invalidated",
			})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("login",       claims.Login)
		c.Set("hospital_id", claims.HospitalID)
		c.Set("hospital",    claims.Hospital)
		c.Set("token",       tokenString)
		c.Set("exp",         claims.ExpiresAt.Time)
		c.Next()
	}
}
