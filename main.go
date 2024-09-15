package main

import (
    "fmt"
    "os" // Add this line
    "os/exec"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var db *gorm.DB

// Define the History and Workspace models
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

// Initialize the database and run the migrations
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

// ExecuteCurlCommand runs a curl command using the system's curl utility
func ExecuteCurlCommand(curlCommand string) (string, error) {
    cmdArgs := strings.Split(curlCommand, " ")
    cmd := exec.Command("curl", cmdArgs...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}

// API handler to execute a curl command and save it in history
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

    // Save the executed command and response to the history
    db.Create(&History{Command: json.Command, Response: response, Timestamp: time.Now()})

    c.JSON(200, gin.H{"response": response})
}

// API handler to retrieve the command execution history
func getHistory(c *gin.Context) {
    var history []History
    db.Order("timestamp desc").Find(&history)
    c.JSON(200, history)
}

// API handler to create a new workspace
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

// API handler to retrieve all workspaces
func getWorkspaces(c *gin.Context) {
    var workspaces []Workspace
    db.Find(&workspaces)
    c.JSON(200, workspaces)
}

// Main function to start the server and define routes
func main() {
    // Initialize the database
    initDB()

    // Create a Gin router
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.Header("Content-Type", "text/html")
        c.String(200, `
        <html>
        <head>
        <title>Brutalist Websites – NOW CHILD</title>

        <style>
        body {
        font-family: monospace;
        background-color: #eee; 
        margin: 40px;
        }
        .box {
        width:400px;
        height: 250px;
        float: left;
        padding: 0px 40px 80px 0px;
        }
        .screenshot {
        width:400px; 
        }
        </style>

        </head>
        <body>

        <div class="box">
        <p>Konata (baguswijaksono)<br/>
        <br/>
        A Lightweight API client testing tool built with Go and a web-based UI. 
        It allows users to execute curl commands via the browser, 
        making it a simple alternative to tools like Postman. <br/>
        <br/>
        <a href="/static/">try now</a><br>
        <a href="https://github.com/baguswijaksono/konata">repo</a>
        </p>
        </div>
        </body>
        </html>
        `)
    })

    r.Static("/static", "./static")

    // Define API routes
    r.POST("/execute", executeCurl)
    r.GET("/history", getHistory)
    r.POST("/workspace", createWorkspace)
    r.GET("/workspaces", getWorkspaces)

    // Start the server
    fmt.Println("Server running at http://localhost:8080")
    r.Run() // listen and serve on 0.0.0.0:8080
}
