package main

import (
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Workspace struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:text"`
	Config    string    `gorm:"type:text"`
	Histories []History `gorm:"foreignKey:WorkspaceID"`
}

type History struct {
	ID          uint      `gorm:"primaryKey"`
	Command     string    `gorm:"type:text"`
	Response    string    `gorm:"type:text"`
	Timestamp   time.Time `gorm:"autoCreateTime"`
	WorkspaceID uint      `gorm:"index;not null;type:bigint unsigned"`
}

func initDB() {
	createDBCmd := exec.Command("mysql", "-u", "phpmyadmin", "-pyour_password", "-e", "CREATE DATABASE IF NOT EXISTS konata")
	if err := createDBCmd.Run(); err != nil {
		panic("failed to create database")
	}

	dsn := "phpmyadmin:your_password@tcp(127.0.0.1:3306)/konata?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	if err := db.AutoMigrate(&Workspace{}, &History{}); err != nil {
		panic("failed to migrate database")
	}
}

func ExecuteCurlCommand(curlCommand string) (string, error) {
	cmdArgs := strings.Split(curlCommand, " ")
	cmd := exec.Command("curl", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func executeCurl(c *gin.Context) {
	var json struct {
		Command   string `json:"command"`
		Workspace string `json:"workspace"` // Expect workspace name from the request
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Find the workspace by name
	var workspace Workspace
	if err := db.Where("name = ?", json.Workspace).First(&workspace).Error; err != nil {
		c.JSON(400, gin.H{"error": "Workspace not found"})
		return
	}

	// Execute the curl command
	response, err := ExecuteCurlCommand(json.Command)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if a similar command exists in the same workspace
	var existingHistory History
	if err := db.Where("command = ? AND workspace_id = ?", json.Command, workspace.ID).First(&existingHistory).Error; err == nil {
		c.JSON(200, gin.H{"response": response})
		return
	}

	// Save the new command history
	db.Create(&History{
		Command:     json.Command,
		Response:    response,
		Timestamp:   time.Now(),
		WorkspaceID: workspace.ID, // Associate the workspace
	})

	c.JSON(200, gin.H{"response": response})
}

func getHistory(c *gin.Context) {
	var history []History
	db.Order("timestamp desc").Find(&history)
	c.JSON(200, history)
}

func createWorkspace(c *gin.Context) {
	var json struct {
		Name   string `json:"name"`
		Config string `json:"config"`
	}
	if err := c.BindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	db.Create(&Workspace{Name: json.Name, Config: json.Config})
	c.JSON(200, gin.H{"message": "Workspace created"})
}

func getWorkspaces(c *gin.Context) {
	var workspaces []Workspace
	db.Find(&workspaces)
	c.JSON(200, workspaces)
}

func main() {
	initDB()
	r := gin.Default()
	r.StaticFile("/", "./static/index.html")

	r.StaticFile("/cli", "./static/cli/index.html")
	r.StaticFile("/gui", "./static/gui/index.html")

	r.POST("/execute", executeCurl)
	r.GET("/history", getHistory)
	r.POST("/workspace", createWorkspace)
	r.GET("/workspaces", getWorkspaces)
	r.Run(":8081")
}
