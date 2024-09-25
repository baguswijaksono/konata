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
		Workspace string `json:"workspace"`
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var workspace Workspace
	if err := db.Where("name = ?", json.Workspace).First(&workspace).Error; err != nil {
		c.JSON(400, gin.H{"error": "Workspace not found"})
		return
	}

	response, err := ExecuteCurlCommand(json.Command)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var existingHistory History
	if err := db.Where("command = ? AND workspace_id = ?", json.Command, workspace.ID).First(&existingHistory).Error; err == nil {
		c.JSON(200, gin.H{"response": response})
		return
	}

	db.Create(&History{
		Command:     json.Command,
		Response:    response,
		Timestamp:   time.Now(),
		WorkspaceID: workspace.ID,
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

// Edit workspace feature
func editWorkspace(c *gin.Context) {
	var json struct {
		Name   string `json:"name"`
		Config string `json:"config"`
	}
	id := c.Param("id")

	if err := c.BindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var workspace Workspace
	if err := db.First(&workspace, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Workspace not found"})
		return
	}

	workspace.Name = json.Name
	workspace.Config = json.Config
	db.Save(&workspace)

	c.JSON(200, gin.H{"message": "Workspace updated"})
}

// Delete workspace feature
func deleteWorkspace(c *gin.Context) {
	id := c.Param("id")

	var workspace Workspace
	if err := db.First(&workspace, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Workspace not found"})
		return
	}

	// Delete associated histories
	db.Where("workspace_id = ?", workspace.ID).Delete(&History{})
	// Delete the workspace
	db.Delete(&workspace)

	c.JSON(200, gin.H{"message": "Workspace deleted"})
}

func getWorkspace(c *gin.Context) {
	id := c.Param("id") // Get the workspace ID from the URL parameters
	var workspace Workspace
	if err := db.First(&workspace, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Workspace not found"})
		return
	}
	c.JSON(200, workspace)
}

func main() {
	initDB()
	r := gin.Default()

	staticFiles := map[string]string{
		"/":       "./static/index.html",
		"/cli":    "./static/html/cli_index.html",
		"/gui":    "./static/html/gui_index.html",
		"/style":  "./static/css/style.css",
		"/cli_js": "./static/js/cli_index.js",
		"/works":  "./static/html/workspace_index.html",

		"/http": "./static/html/gui_http.html",
		"/ws":   "./static/html/gui_ws.html",
		"/mqtt": "./static/html/gui_mqqt.html",
		"/tcp":  "./static/html/gui_tcp.html",
		"/ftp":  "./static/html/gui_ftp.html",
		"/smtp": "./static/html/gui_smtp.html",
		"/pop3": "./static/html/gui_pop3.html",
		"/imap": "./static/html/gui_imap.html",
	}

	for route, file := range staticFiles {
		r.StaticFile(route, file)
	}

	r.POST("/execute", executeCurl)
	r.GET("/history", getHistory)
	r.POST("/workspace", createWorkspace)
	r.GET("/workspaces", getWorkspaces)
	r.GET("/workspace/:id", getWorkspace)       // Get workspace route
	r.PUT("/workspace/:id", editWorkspace)      // Edit workspace route
	r.DELETE("/workspace/:id", deleteWorkspace) // Delete workspace route

	r.Run()
}
