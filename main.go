package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Task struct represents the structure of a task
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

var db *sql.DB

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Define API endpoints
	// r.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	r.POST("/tasks", createTask)
	r.GET("/tasks/:id", getTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.GET("/tasks", listTasks)

	// Run the server
	r.Run(":8080")
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./task.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Create tasks table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
                            id INTEGER PRIMARY KEY ,
                            title TEXT,
                            description TEXT,
                            due_date TEXT,
                            status TEXT
                        )`)
	if err != nil {
		log.Fatal("Error creating tasks table:", err)
	}
}

// Create a new task
func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
     
	// Insert task into database
	result, err := db.Exec("INSERT INTO tasks (id,title, description, due_date, status) VALUES (? ,? , ?, ?, ?)",
		task.ID, task.Title, task.Description, task.DueDate, task.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Get the ID of the newly inserted task
	id, _ := result.LastInsertId()
	task.ID = int(id)

	c.JSON(http.StatusCreated, task)
}

// Retrieve a task by ID
func getTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var task Task
	err := db.QueryRow("SELECT id, title, description, due_date, status FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
	// c.JSON(200, gin.H{
	// 	"message": "pong",
	// })
}

// Update a task by ID
func updateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Give description title status"})
		return
	}
	exists, err := idExists(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if ID exists"})
        return
    }

    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
        return
    }
	// Update task in the database
	_,e := db.Exec("UPDATE tasks SET title=?, description=?, due_date=?, status=? WHERE id=?",
		task.Title, task.Description, task.DueDate, task.Status, id)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	
	task.ID = id
	c.JSON(http.StatusOK, task)
   
}

// Delete a task by ID
func deleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	exists, err := idExists(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if ID exists"})
        return
    }

    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
        return
    }

	// Delete task from the database
	_, e := db.Exec("DELETE FROM tasks WHERE id=?", id)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// List all tasks
func listTasks(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, description, due_date, status FROM tasks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}
//check that id exist in table or not
func idExists(id int) (bool, error) {
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}