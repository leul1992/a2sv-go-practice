Task Manager API

REST API for managing tasks (CRUD) using Go and Gin.

MongoDB Integration

This version integrates MongoDB for persistent data storage using the MongoDB Go Driver.

Prerequisites:
- Install MongoDB and start the server locally (default: mongodb://localhost:27017), or use a cloud provider like MongoDB Atlas.
- Database: taskdb
- Collection: tasks
- have .env at the root of the folder with a varable "MONGODB_URI" which is uri of local or cloud mongodb.

Setup Instructions:
1. Clone the repository and navigate to the task_manager folder.
2. Initialize the Go module: go mod init task_manager
3. Install dependencies:
   - go get -u github.com/gin-gonic/gin
   - go get go.mongodb.org/mongo-driver/mongo
4. Start your MongoDB instance.
5. Run the application: go run main.go

The application will connect to MongoDB on startup. Verify connection by checking logs.

To validate data:
- Use MongoDB Compass or mongo shell to query the 'tasks' collection in 'taskdb' database.
- Example query: db.tasks.find({})

API Endpoints (unchanged from previous version):

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

Testing:
- Use Postman or curl to test CRUD operations.
- After operations, verify data persistence by restarting the server and querying endpoints or MongoDB directly.
- Error handling includes 404 for not found, 400 for invalid input, 500 for database errors.