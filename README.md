# Task Management API using Gogin Framework

This repository contains the implementation of a Task Management API using Gogin Framework. The API allows users to perform CRUD operations on tasks stored in a SQLite database.

## Table of Contents
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
  - [Create a New Task](#create-a-new-task)
  - [Retrieve a Task](#retrieve-a-task)
  - [Update a Task](#update-a-task)
  - [Delete a Task](#delete-a-task)
  - [List All Tasks](#list-all-tasks)

## Database Schema
The SQLite database schema for storing tasks includes the following fields:
- ID: The unique identifier for each task.
- Title: The title or name of the task.
- Description: A brief description of the task.
- Due Date: The deadline or due date for the task.
- Status: The status of the task (e.g., "pending," "completed," "in progress," etc.).

## API Endpoints

### Create a New Task
- **Endpoint:** `POST /tasks`
- **Description:** Accepts a JSON payload containing the task details (title, description, due date). Generates a unique ID for the task and stores it in the database. Returns the created task with the assigned ID.

### Retrieve a Task
- **Endpoint:** `GET /tasks/{id}`
- **Description:** Accepts a task ID as a parameter. Retrieves the corresponding task from the database. Returns the task details if found, or an appropriate error message if not found.

### Update a Task
- **Endpoint:** `PUT /tasks/{id}`
- **Description:** Accepts a task ID as a parameter. Accepts a JSON payload containing the updated task details (title, description, due date). Updates the corresponding task in the database. Returns the updated task if successful, or an appropriate error message if not found.

### Delete a Task
- **Endpoint:** `DELETE /tasks/{id}`
- **Description:** Accepts a task ID as a parameter. Deletes the corresponding task from the database. Returns a success message if the deletion is successful, or an appropriate error message if not found.

### List All Tasks
- **Endpoint:** `GET /tasks`
- **Description:** Retrieves all tasks from the database. Returns a list of tasks, including their details (title, description, due date).

## Usage
1. Clone this repository.
2. Install the necessary dependencies.
3. Run the server.
4. Access the API endpoints using appropriate HTTP requests.

## Technologies Used
- Gogin Framework
- SQLite
- Go programming language
