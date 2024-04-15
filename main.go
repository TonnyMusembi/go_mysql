package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Define a global variable for the database connection
var db *sql.DB

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Info("This is an informational log message")

	// Connect to the MySQL database

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Retrieve database connection parameters from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the MySQL connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Connect to the MySQL database
	// var err error
	db, err = sql.Open("mysql", connectionString)
	log.Println("connected")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()
	router := gin.Default()

	// Define routes
	router.POST("/users", addUsers)
	router.GET("/users", getUsers)

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
		Id   int    `json:"id"`
		Name string `json:"name"`
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

	_, err = stmt.Exec(user.Id, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}
