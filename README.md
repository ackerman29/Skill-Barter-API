
#  Skill Barter Platform (Backend)

## LIVE : https://skill-barter-api.onrender.com/

This is the backend API for a **Skill Barter Platform for Students**, where users can offer and request skills, match with others, and exchange learning sessions. Built using **Go (Gin)** and **MongoDB**.


## DEMO WALTHROUGH:
https://drive.google.com/file/d/1ghLobAz04YG58X3ccJamHeiw0FZCTkl2/view?usp=drive_link
---
---



##  Tech Stack

- **Backend:** Go, Gin
- **Database:** MongoDB Atlas
- **Authentication:** JWT + bcrypt
- **API Format:** REST (JSON)
- **Middleware:** CORS, JWT Auth
- **WebSockets:** Gorilla Websockets for Real-time communication for requests/responses
- **Hosting:** Render


--------------------------------------------------------------------------------------------------------------


##  Setup Instructions

1. **Clone the repository**

```bash
git clone https://github.com/your-username/skill-barter-backend.git
cd skill-barter-backend
2. **Install dependencies**

```bash
go mod tidy
```

3. **Start MongoDB** (if not already running)

4. **Run the server**

```bash
go run main.go
```

The server will start on `http://localhost:8000`

--------------------------------------------------------------------------------------------------------------

##  Authentication Routes

### POST `/auth/signup`

Create a new user.

```json
Request Body:
{
  "email": "test@example.com",
  "password": "yourpassword"
}
```

--------------------------------------------------------------------------------------------------------------

### POST `/auth/login`

Login with credentials.

```json
Request Body:
{
  "email": "test@example.com",
  "password": "yourpassword"
}
```

**Response:**
```json
{
  "token": "JWT-TOKEN"
}
```

--------------------------------------------------------------------------------------------------------------

##  Protected Routes

Add `Authorization: Bearer <JWT>` header in requests.

### GET `/auth/myprofile`

Fetch user profile.

--------------------------------------------------------------------------------------------------------------

### PUT `/auth/myprofile`

Update profile.

```json
{
  "skillsHave": ["Go", "React"],
  "skillsWant": ["Python", "Docker"],
  "availableDays": ["Saturday", "Sunday"]
}
```

--------------------------------------------------------------------------------------------------------------

### GET `/match`

Returns a list of matched users based on skills.

--------------------------------------------------------------------------------------------------------------

##  Protected Test Route

### GET `/protected`

Just to test JWT middleware.

--------------------------------------------------------------------------------------------------------------

##  CORS Setup

If using frontend on `localhost:5173`, CORS is enabled using:

```go
import "github.com/gin-contrib/cors"

r.Use(cors.Default())
```

--------------------------------------------------------------------------------------------------------------

##  Testing with Postman

You can test the APIs using tools like **Postman** or **curl**:

1. Signup → Get JWT token
2. Login → Get JWT token
3. Use JWT in protected routes
4. Test `/auth/myprofile` and `/match`

--------------------------------------------------------------------------------------------------------------
