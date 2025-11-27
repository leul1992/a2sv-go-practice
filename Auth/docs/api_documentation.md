Task Manager API

REST API for managing tasks (CRUD) using Go and Gin.

MongoDB Integration

This version integrates MongoDB for persistent data storage using the MongoDB Go Driver.

Prerequisites:
- Install MongoDB and start the server locally (default: mongodb://localhost:27017), or use a cloud provider like MongoDB Atlas.
- Database: taskdb
- Collections: tasks and users
- have .env at the root of the folder with these variables:

  MONGODB_URI=mongodb://localhost:27017
  JWT_SECRET=your_very_strong_secret_key_here

Setup Instructions:
1. Clone the repository and navigate to the task_manager folder.
2. Initialize the Go module: go mod init task_manager
3. Install dependencies:
   - go get -u github.com/gin-gonic/gin
   - go get go.mongodb.org/mongo-driver/mongo
   - go get github.com/joho/godotenv
   - go get golang.org/x/crypto/bcrypt
   - go get github.com/golang-jwt/jwt/v5
4. Start your MongoDB instance.
5. Run the application: go run main.go

The application will connect to MongoDB on startup. Verify connection by checking logs.

To validate data:
- Use MongoDB Compass or mongo shell to query the 'tasks' and 'users' collections in 'taskdb' database.
- Example query: db.tasks.find({}) or db.users.find({})

Authentication (JWT + Role-based)

New Endpoints:

POST — Register User
URL: http://localhost:8080/register

Body (raw JSON):
{
  "username": "Leul",
  "password": "Leul123"
}

POST — Login (User or Admin)
URL: http://localhost:8080/login

Body (raw JSON):
{
  "username": "Leul",
  "password": "Leul123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "user",
  "user_id": "66f8a1b2c123456789abcde"
}

POST — Promote User to Admin (Admin only)
URL: http://localhost:8080/promote/:id

Example: http://localhost:8080/promote/66f8a1b2c123456789abcde

Header:
Authorization: Bearer <admin-jwt-token>

All Task Endpoints Now Require Authentication

Add this header to every task request:
Authorization: Bearer <token-from-login>

API Endpoints (now protected):

GET — Get Tasks
URL: http://localhost:8080/tasks
Header: Authorization: Bearer <token>

GET — Get Task
URL: http://localhost:8080/tasks/1
Header: Authorization: Bearer <token>

POST — Create Task
URL: http://localhost:8080/tasks
Header: Authorization: Bearer <token>

Body (raw JSON):
{
  "title": "New Task",
  "description": "Learning Go with Gin",
  "due_date": "2025-12-01",
  "status": "pending"
}

PUT — Update Task
URL: http://localhost:8080/tasks/1
Header: Authorization: Bearer <token>

Body (raw JSON):
{
  "title": "Updated Task",
  "description": "Updated description",
  "due_date": "2025-10-10",
  "status": "completed"
}

DELETE — Delete Task
URL: http://localhost:8080/tasks/1
Header: Authorization: Bearer <token>

Testing:
- Use Postman or curl to test CRUD operations.
- First register → login → copy token → set Authorization header as Bearer Token in Postman.
- After operations, verify data persistence by restarting the server and querying endpoints or MongoDB directly.
- Error handling includes 404 for not found, 400 for invalid input, 500 for database errors.
- Also: 401 Unauthorized (missing/invalid token), 403 Forbidden (not admin when needed).
