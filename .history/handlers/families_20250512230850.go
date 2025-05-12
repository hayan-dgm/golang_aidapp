package handlers

import (
	"aidapp_api_golang/db"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GetFamilies - Matches Python's paginated response exactly
func GetFamilies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	// Get column names for mapping
	rows, err := db.DB.Query("SELECT * FROM families LIMIT ? OFFSET ?", perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var families []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data parsing error"})
			return
		}

		family := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				family[col] = string(b)
			} else {
				family[col] = val
			}
		}
		families = append(families, family)
	}

	var total int
	db.DB.QueryRow("SELECT COUNT(*) FROM families").Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"page":     page,
		"per_page": perPage,
		"total":    total,
		"families": families,
	})
}

// GetFamily - Matches Python's dynamic column handling
// func GetFamily(c *gin.Context) {
// 	id := c.Param("id")

// 	row := db.DB.QueryRow("SELECT * FROM families WHERE id = ?", id)

// 	columns, _ := row.Columns()
// 	values := make([]interface{}, len(columns))
// 	pointers := make([]interface{}, len(columns))
// 	for i := range values {
// 		pointers[i] = &values[i]
// 	}

// 	if err := row.Scan(pointers...); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			c.JSON(http.StatusNotFound, gin.H{"message": "Family not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		}
// 		return
// 	}

// 	family := make(map[string]interface{})
// 	for i, col := range columns {
// 		val := values[i]
// 		if b, ok := val.([]byte); ok {
// 			family[col] = string(b)
// 		} else {
// 			family[col] = val
// 		}
// 	}

// 	c.JSON(http.StatusOK, family)
// }

func GetFamily(c *gin.Context) {
	id := c.Param("id")

	var family struct {
		ID            int    `json:"id"`
		FullName      string `json:"fullName"`
		NationalID    string `json:"nationalID"`
		FamilyBookID  string `json:"familyBookID"`
		PhoneNumber   string `json:"phoneNumber"`
		FamilyMembers int    `json:"familyMembers"`
		Children      int    `json:"children"`
		Babies        int    `json:"babies"`
		Adults        int    `json:"adults"`
		Milk          int    `json:"milk"`
		Diapers       int    `json:"diapers"`
		Basket        int    `json:"basket"`
		Clothing      int    `json:"clothing"`
		Drugs         int    `json:"drugs"`
		Other         string `json:"other"`
		Taken         bool   `json:"taken"`
	}

	err := db.DB.QueryRow(`
        SELECT id, fullName, nationalID, familyBookID, phoneNumber,
               familyMembers, children, babies, adults,
               milk, diapers, basket, clothing, drugs, other, taken
        FROM families WHERE id = ?`, id).Scan(
		&family.ID, &family.FullName, &family.NationalID, &family.FamilyBookID, &family.PhoneNumber,
		&family.FamilyMembers, &family.Children, &family.Babies, &family.Adults,
		&family.Milk, &family.Diapers, &family.Basket, &family.Clothing, &family.Drugs, &family.Other, &family.Taken)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Family not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, family)
}

// AddFamily - Matches Python's admin check and field validation
func AddFamily(c *gin.Context) {
	claims := c.MustGet("jwt_claims").(jwt.MapClaims)
	user := claims["sub"].(map[string]interface{})

	if !user["isAdmin"].(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req struct {
		FullName      string `json:"fullName" binding:"required"`
		NationalID    string `json:"nationalID" binding:"required"`
		FamilyBookID  string `json:"familyBookID" binding:"required"`
		PhoneNumber   string `json:"phoneNumber" binding:"required"`
		FamilyMembers int    `json:"familyMembers" binding:"required"`
		Children      int    `json:"children" binding:"required"`
		Babies        int    `json:"babies" binding:"required"`
		Adults        int    `json:"adults" binding:"required"`
		Milk          int    `json:"milk"`
		Diapers       int    `json:"diapers"`
		Basket        int    `json:"basket"`
		Clothing      int    `json:"clothing"`
		Drugs         int    `json:"drugs"`
		Other         string `json:"other"`
		Taken         bool   `json:"taken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO families (
		fullName, nationalID, familyBookID, phoneNumber,
		familyMembers, children, babies, adults,
		milk, diapers, basket, clothing, drugs, other, taken
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.FullName, req.NationalID, req.FamilyBookID, req.PhoneNumber,
		req.FamilyMembers, req.Children, req.Babies, req.Adults,
		req.Milk, req.Diapers, req.Basket, req.Clothing, req.Drugs, req.Other, req.Taken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create family"})
		return
	}

	id, _ := res.LastInsertId()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	_, err = tx.Exec(`
		INSERT INTO logs (familyID, userID, changeDescription, timestamp)
		VALUES (?, ?, ?, ?)`,
		id, int(user["id"].(float64)), "Family created", timestamp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log action"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Family added successfully",
		"familyId": id,
	})
}

// UpdateProducts - Matches Python's COALESCE behavior and WebSocket notifications
func UpdateProducts(c *gin.Context) {
	id := c.Param("id")
	claims := c.MustGet("jwt_claims").(jwt.MapClaims)
	user := claims["sub"].(map[string]interface{})

	var req struct {
		Milk     *int    `json:"milk"`
		Diapers  *int    `json:"diapers"`
		Basket   *int    `json:"basket"`
		Clothing *int    `json:"clothing"`
		Drugs    *int    `json:"drugs"`
		Other    *string `json:"other"`
		Taken    *bool   `json:"taken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer tx.Rollback()

	// Get current values for comparison
	var current struct {
		Milk     int
		Diapers  int
		Basket   int
		Clothing int
		Drugs    int
		Other    string
		Taken    bool
	}
	err = tx.QueryRow(`
		SELECT milk, diapers, basket, clothing, drugs, other, taken
		FROM families WHERE id = ?`, id).Scan(
		&current.Milk, &current.Diapers, &current.Basket,
		&current.Clothing, &current.Drugs, &current.Other, &current.Taken)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Family not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Build dynamic update
	changes := make(map[string]interface{})
	query := "UPDATE families SET "
	args := []interface{}{}
	updates := []string{}

	addUpdate := func(field string, value interface{}, currentVal interface{}) {
		if value != nil && value != currentVal {
			updates = append(updates, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
			changes[field] = value
		}
	}

	addUpdate("milk", req.Milk, current.Milk)
	addUpdate("diapers", req.Diapers, current.Diapers)
	addUpdate("basket", req.Basket, current.Basket)
	addUpdate("clothing", req.Clothing, current.Clothing)
	addUpdate("drugs", req.Drugs, current.Drugs)
	addUpdate("other", req.Other, current.Other)
	addUpdate("taken", req.Taken, current.Taken)

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query += strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err = tx.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	// Log the change
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = tx.Exec(`
		INSERT INTO logs (familyID, userID, changeDescription, timestamp)
		VALUES (?, ?, ?, ?)`,
		id, int(user["id"].(float64)), "Products updated", timestamp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log action"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	// WebSocket notification
	if len(changes) > 0 {
		hub.Broadcast <- Message{
			Room: "family_" + id,
			Payload: gin.H{
				"family_id": id,
				"changes":   changes,
				"updated_by": gin.H{
					"user_id":  user["id"],
					"username": user["username"],
				},
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// ActiveSessions - Matches Python's implementation
func GetActiveSessions(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT active_sessions.user_id, users.username, active_sessions.login_time 
		FROM active_sessions 
		JOIN users ON active_sessions.user_id = users.id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var sessions []gin.H
	for rows.Next() {
		var userID int
		var username string
		var loginTime time.Time
		if err := rows.Scan(&userID, &username, &loginTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data parsing error"})
			return
		}
		sessions = append(sessions, gin.H{
			"user_id":    userID,
			"username":   username,
			"login_time": loginTime.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, sessions)
}

func ClearActiveSessions(c *gin.Context) {
	if _, err := db.DB.Exec("DELETE FROM active_sessions"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All active sessions cleared"})
}
