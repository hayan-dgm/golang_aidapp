// package middleware

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v4"
// )

// type JWTConfig struct {
// 	SecretKey      string
// 	AccessTokenExp time.Duration
// }

// var (
// 	Config = JWTConfig{
// 		SecretKey:      getEnv("JWT_SECRET_KEY", "default-secret-key"),
// 		AccessTokenExp: time.Hour * 1,
// 	}
// 	DB *sql.DB
// )

// func Initialize(db *sql.DB) {
// 	DB = db
// }

// func JWTMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		if tokenString == authHeader {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Bearer token required"})
// 			return
// 		}

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return []byte(Config.SecretKey), nil
// 		})

// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			// Check if token is revoked
// 			jti, ok := claims["jti"].(string)
// 			if !ok {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
// 				return
// 			}

// 			var revoked bool
// 			err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE jti = ?)", jti).Scan(&revoked)
// 			if err != nil {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
// 				return
// 			}

// 			if revoked {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token has been revoked"})
// 				return
// 			}

// 			// Check session timestamp
// 			sub, ok := claims["sub"].(map[string]interface{})
// 			if !ok {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
// 				return
// 			}

// 			userID, ok := sub["id"].(float64)
// 			if !ok {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
// 				return
// 			}

// 			var dbLoginTime time.Time
// 			err = DB.QueryRow("SELECT login_time FROM active_sessions WHERE user_id = ?", int(userID)).Scan(&dbLoginTime)
// 			if err != nil && err != sql.ErrNoRows {
// 				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
// 				return
// 			}

// 			if !dbLoginTime.IsZero() {
// 				iat, ok := claims["iat"].(float64)
// 				if !ok {
// 					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
// 					return
// 				}

// 				tokenLoginTime := time.Unix(int64(iat), 0)
// 				if dbLoginTime.After(tokenLoginTime.Add(time.Second)) {
// 					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token has been revoked"})
// 					return
// 				}
// 			}

// 			// Store claims in context
// 			c.Set("jwt_claims", claims)
// 			c.Next()
// 		} else {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
// 		}
// 	}
// }

// func GenerateToken(userID int, username string, isAdmin bool) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": map[string]interface{}{
// 			"id":       userID,
// 			"username": username,
// 			"isAdmin":  isAdmin,
// 		},
// 		"iat": time.Now().Unix(),
// 		"exp": time.Now().Add(Config.AccessTokenExp).Unix(),
// 		"jti": generateJTI(),
// 	})

// 	return token.SignedString([]byte(Config.SecretKey))
// }

// func generateJTI() string {
// 	return fmt.Sprintf("%d", time.Now().UnixNano())
// }

// func getEnv(key, defaultValue string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultValue
// }
// =================================================

package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTConfig struct {
	SecretKey      string
	AccessTokenExp time.Duration
}

var (
	Config = JWTConfig{
		SecretKey:      getEnv("JWT_SECRET_KEY", "default-secret-key"),
		AccessTokenExp: time.Hour * 1,
	}
	DB *sql.DB
)

func Initialize(db *sql.DB) {
	DB = db
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Bearer token required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(Config.SecretKey), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check if token is revoked
			jti, ok := claims["jti"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
				return
			}

			var revoked bool
			err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE jti = ?)", jti).Scan(&revoked)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
				return
			}

			if revoked {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token has been revoked"})
				return
			}

			// Check session timestamp
			sub, ok := claims["sub"].(map[string]interface{})
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
				return
			}

			userID, ok := sub["id"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
				return
			}

			var dbLoginTime time.Time
			err = DB.QueryRow("SELECT login_time FROM active_sessions WHERE user_id = ?", int(userID)).Scan(&dbLoginTime)
			if err != nil && err != sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
				return
			}

			if !dbLoginTime.IsZero() {
				iat, ok := claims["iat"].(float64)
				if !ok {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
					return
				}

				tokenLoginTime := time.Unix(int64(iat), 0)
				if dbLoginTime.After(tokenLoginTime.Add(time.Second)) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token has been revoked"})
					return
				}
			}

			// Store claims in context
			c.Set("jwt_claims", claims)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		}
	}
}

// func GenerateToken(userID int, username string, isAdmin bool) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": map[string]interface{}{
// 			"id":       userID,
// 			"username": username,
// 			"isAdmin":  isAdmin,
// 		},
// 		"iat": time.Now().Unix(),
// 		"exp": time.Now().Add(Config.AccessTokenExp).Unix(),
// 		"jti": generateJTI(),
// 	})

// return token.SignedString([]byte(Config.SecretKey))}
func generateJTI() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func GenerateToken(userID int, username string, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": map[string]interface{}{
			"id":       userID,
			"username": username,
			"isAdmin":  isAdmin,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(Config.AccessTokenExp).Unix(),
		"jti": generateJTI(),
	})

	return token.SignedString([]byte(Config.SecretKey))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
