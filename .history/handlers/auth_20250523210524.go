// package handlers

// import (
// 	"database/sql"
// 	"net/http"
// 	"golang_api_aidapp/middleware"
// 	"yourproject/models"

// 	"github.com/gin-gonic/gin"
// 	"golang.org/x/crypto/bcrypt"
// )

// func Login(c *gin.Context) {
// 	var creds struct {
// 		Username string `json:"username" binding:"required"`
// 		Password string `json:"password" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&creds); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var user models.User
// 	err := db.DB.QueryRow(`
// 		SELECT id, username, password, isAdmin
// 		FROM users WHERE username = ?`, creds.Username).Scan(
// 		&user.ID, &user.Username, &user.Password, &user.IsAdmin)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		}
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
// 		return
// 	}

// 	token, err := middleware.GenerateToken(user.ID, user.Username, user.IsAdmin)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	// Store session
// 	_, err = db.DB.Exec(`
// 		INSERT INTO active_sessions (user_id, access_token, login_time)
// 		VALUES (?, ?, ?)`,
// 		user.ID, token, time.Now().UTC())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"access_token": token,
// 		"message":      "Login successful",
// 	})
// }

// func Logout(c *gin.Context) {
// 	claims := c.MustGet("jwt_claims").(jwt.MapClaims)
// 	jti := claims["jti"].(string)
// 	userID := claims["sub"].(map[string]interface{})["id"].(float64)

// 	// Delete session and revoke token
// 	_, err := db.DB.Exec(`
// 		DELETE FROM active_sessions WHERE user_id = ?;
// 		INSERT INTO revoked_tokens (jti) VALUES (?)`,
// 		int(userID), jti)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
// }

package handlers

import (
	"aidapp_api_golang/db"
	"aidapp_api_golang/middleware"
	"aidapp_api_golang/models"
	"database/sql"
	"net/http"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API"})
}

func Login(c *gin.Context) {
	var creds struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := db.DB.QueryRow(`
		SELECT id, username, password, isAdmin 
		FROM users WHERE username = ?`, creds.Username).Scan(
		&user.ID, &user.Username, &user.Password, &user.IsAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Store session
	_, err = db.DB.Exec(`
		INSERT INTO active_sessions (user_id, access_token, login_time) 
		VALUES (?, ?, ?)`,
		user.ID, token, time.Now().UTC())
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
	// 	return
	// }
	if err != nil {
		log.Printf("Failed to create session: %v", err) // log the real error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"message":      "Login successful",
	})
}

func Logout(c *gin.Context) {
	claims := c.MustGet("jwt_claims").(jwt.MapClaims)
	jti := claims["jti"].(string)
	userID := claims["sub"].(map[string]interface{})["id"].(float64)

	// Delete session and revoke token
	_, err := db.DB.Exec(`
		DELETE FROM active_sessions WHERE user_id = ?;
		INSERT INTO revoked_tokens (jti) VALUES (?)`,
		int(userID), jti)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
