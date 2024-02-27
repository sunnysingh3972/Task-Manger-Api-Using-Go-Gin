package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/sunnysingh3972/Task-Manger-Api-Using-Go-Gin/database"
	"github.com/sunnysingh3972/Task-Manger-Api-Using-Go-Gin/handlers"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	taskHandler := handlers.NewTaskHandler(db)

	r.POST("/tasks", taskHandler.CreateTask)
	r.GET("/tasks/:id", taskHandler.GetTask)
	r.PUT("/tasks/:id", taskHandler.UpdateTask)
	r.DELETE("/tasks/:id", taskHandler.DeleteTask)
	r.GET("/tasks", taskHandler.ListTasks)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
