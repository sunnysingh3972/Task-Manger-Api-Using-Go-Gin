package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/sunnysingh3972/Task-Manger-Api-Using-Go-Gin/models"
)

type TaskHandler struct {
	DB *sql.DB
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.DB.Exec("INSERT INTO tasks (title, description, due_date, status) VALUES (?, ?, ?, ?)",
		task.Title, task.Description, task.DueDate, task.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	id, _ := result.LastInsertId()
	task.ID = int(id)

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	err = h.DB.QueryRow("SELECT id, title, description, due_date, status FROM tasks WHERE id = ?", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Give description title status"})
		return
	}
	exists, err := h.idExists(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if ID exists"})
        return
    }

    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
        return
    }
	// Update task in the database
	_,e := h.DB.Exec("UPDATE tasks SET title=?, description=?, due_date=?, status=? WHERE id=?",
		task.Title, task.Description, task.DueDate, task.Status, id)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	
	task.ID = id
	c.JSON(http.StatusOK, task)
   
}

// Delete a task by ID
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	exists, err := h.idExists(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if ID exists"})
        return
    }

    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
        return
    }

	// Delete task from the database
	_, e := h.DB.Exec("DELETE FROM tasks WHERE id=?", id)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// List all tasks
func (h *TaskHandler) ListTasks(c *gin.Context) {
	rows, err := h.DB.Query("SELECT id, title, description, due_date, status FROM tasks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}

//check that id exist in table or not
func (h *TaskHandler) idExists(id int) (bool, error) {
    var exists bool
    err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}
