Task Manager API

REST API for managing tasks (CRUD) using Go and Gin.

Step By Step Guide:

First Clone the folder

Make sure you cd to task_manager

Open your terminal in the task_manager folder:

  - go mod init task_manager


Then install Gin Framework:

  - go get -u github.com/gin-gonic/gin

From the root project folder:

  - go run main.go


If everything is correct, you will see something like:

  Listening and serving HTTP on :8080


Tasks You can do:

  GET — Get Tasks

  URL:
  http://localhost:8080/tasks

  GET — Get Task

  URL:
  http://localhost:8080/tasks/1

  POST — Create Task

  URL:
  http://localhost:8080/tasks

  Body (raw JSON):
  {
    "title": "New Task",
    "description": "Learning Go with Gin",
    "due_date": "2025-12-01",
    "status": "pending"
  }

  PUT — Update Task

  URL:
  http://localhost:8080/tasks/1

  Body (raw JSON):
  {
    "title": "Updated Task",
    "description": "Updated description",
    "due_date": "2025-10-10",
    "status": "completed"
  }

  DELETE — Delete Task

  URL:
  http://localhost:8080/tasks/1