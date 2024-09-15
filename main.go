package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var db *gorm.DB
type History struct {
    ID        uint      `gorm:"primaryKey"`
    Command   string    `gorm:"type:text"`
    Response  string    `gorm:"type:text"`
    Timestamp time.Time `gorm:"autoCreateTime"`
}

type Workspace struct {
    ID      uint   `gorm:"primaryKey"`
    Name    string `gorm:"type:text"`
    Config  string `gorm:"type:text"`
}

func initDB() {
    if _, err := os.Stat("history.db"); os.IsNotExist(err) {
        fmt.Println("Database file does not exist, creating a new one...")
    } else if err != nil {
        panic("failed to check if database file exists")
    }

    var err error
    db, err = gorm.Open(sqlite.Open("history.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database")
    }
    db.AutoMigrate(&History{}, &Workspace{})
}

func ExecuteCurlCommand(curlCommand string) (string, error) {
    cmdArgs := strings.Split(curlCommand, " ")
    cmd := exec.Command("curl", cmdArgs...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}

func executeCurl(c *gin.Context) {
    var json struct {
        Command string `json:"command"`
    }
    if err := c.BindJSON(&json); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    response, err := ExecuteCurlCommand(json.Command)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    db.Create(&History{Command: json.Command, Response: response, Timestamp: time.Now()})
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
