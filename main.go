package main

import (
	"database/sql"
	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Define a global variable for the database connection
var db *sql.DB

func main() {
    // Connect to the MySQL database
    var err error
    db, err = sql.Open("mysql", "root:tonny@07@tcp(localhost:3306)/movies")
    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }

    // Ping the database to check if the connection is successful
    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database:", err)
    }

    // Create a new Gin router
    router := gin.Default()

    // Define routes
    router.GET("/users", getUsers)
    router.POST("/users", addUsers)
	
    // Start the server
    err = router.Run(":8080")
    if err != nil {
        log.Fatal("Error starting server:", err)
    }
}

func getUsers(c *gin.Context) {
    // Perform a SELECT query to retrieve users from the database
    rows, err := db.Query("SELECT id,  name FROM users")
    if err != nil {
        c.JSON(500, gin.H{"error": "Error querying the database"})
        return
    }
    defer rows.Close()

    // Iterate over the result set and build a slice of users
    var users []gin.H
    for rows.Next() {
        var id int
        var name string
        err := rows.Scan(&id, &name)
        if err != nil {
            c.JSON(500, gin.H{"error": "Error scanning rows"})
            return
        }
        users = append(users, gin.H{"id": id, "username": name})
    }

    // Check for any errors encountered during iteration
    err = rows.Err()
    if err != nil {
        c.JSON(500, gin.H{"error": "Error iterating over rows"})
        return
    }

    // Return the list of users as JSON
    c.JSON(200, users)
}
func addUsers(c *gin.Context) {
    // Parse JSON request body
    var user struct {
        Id int `json:"id"`
        Name    string `json:"name"`
        
    }
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Insert user into the database
    stmt, err := db.Prepare("INSERT INTO users (id, name) VALUES (?, ?)")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing SQL statement"})
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(user.Id, user.Name,)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}
