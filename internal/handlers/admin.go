package handlers

import (
	"net/http"
	"repo-backend/internal/database"
	"repo-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateAdmin(c *gin.Context) {
	var admin models.Admin

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&admin)

	c.JSON(http.StatusOK, gin.H{"message": "Admin created"})
}

func GetAdmins(c *gin.Context) {
	admins := []models.Admin{}

	database.DB.Find(&admins)

	c.JSON(http.StatusOK, admins)
}

func GetAdmin(c *gin.Context) {
	id := c.Param("id")

	var admin models.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	c.JSON(http.StatusOK, admin)
}

func UpdateAdmin(c *gin.Context) {
	id := c.Param("id")

	var admin models.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&admin)

	c.JSON(http.StatusOK, gin.H{"message": "Admin updated"})
}

func DeleteAdmin(c *gin.Context) {
	id := c.Param("id")

	var admin models.Admin
	if err := database.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	database.DB.Delete(&admin)

	c.JSON(http.StatusOK, gin.H{"message": "Admin deleted"})
}
