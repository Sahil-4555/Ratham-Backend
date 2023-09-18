# Ratham-Backend

Ratham Backend is a Go-based web application that manages sessions between students and deans at a university. It provides a set of APIs for user registration, login, session booking, and management.

## Features & Functionality

Ratham Backend offers the following key features and functionalities:

1. Session Management

    - **Session Booking**: Students can view and book available dean sessions. Each dean session is one hour long, scheduled on Thursdays and Fridays at 10 AM every week.
    - **Pending Sessions**: Deans can log in to view a list of pending sessions, including student names and session details.

2. User Management

    - **Student Registration**: New students can register with their university ID and password.
    - **Dean Registration**: New deans can register with their university ID and password.
    - **User Login**: Registered students and deans can log in to the system.
    - **User Logout**: Users can log out to end their sessions securely.

3. Upcoming Free Sessions

    - **Upcoming Sessions**: Deans can log in to see a list of upcoming free sessions, including session start and end times.

4. Authentication
    
    - **JWT Authentication**: User authentication is implemented using JSON Web Tokens (JWT) for secure access to protected routes. To access protected routes, include the JWT token in the request headers as a Bearer token.

## Getting Started

To run this application locally, follow these steps:

1. Clone the repository to your local machine.
```
git clone https://github.com/Sahil-4555/ratham-backend.git
```

2. Navigate to the project directory.
```
cd ratham-backend
```

3. Install the required dependencies.
```
go mod tidy
```

4. Start the server
```
go run main.go
```

**The Server will start on port 8080 by default.**

## API EndPoints

### Students APIs

1. Register:
```
// Register a new student. --POST Method
http://localhost:8080/student/register
```

2. Login:
```
// Authenticate and log in a student. --POST Method
http://localhost:8080/student/login
```

3. Get Student Details:
```
// Retrieve student details. --GET Method
http://localhost:8080/student/user
```

4. Logout:
```
// Log out the student. --POST Method
http://localhost:8080/student/logout
```

### Dean APIs

1. Register:
```
// Register a new dean. --POST Method
http://localhost:8080/dean/register
```

2. Login:
```
// Authenticate and log in a dean. --POST Method
http://localhost:8080/dean/login
```

3. Get Student Details:
```
// Retrieve dean details. --GET Method
http://localhost:8080/dean/user
```

4. Logout:
```
// Log out the dean. --POST Method
http://localhost:8080/dean/logout
```

### Session APIs

1. Add New Session:
```
// Add a new session (deans only) --POST Method
http://localhost:8080/session/addsession
```

2. Get All Free Sessions:
```
// Get a list of free sessions (students only). --GET Method
http://localhost:8080/session/getfreesession
```

3. Book a Session:
```
// Book a session (students only). --POST Method
http://localhost:8080/session/booksession/:sessionid
```

4. Get Upcoming Free Sessions:
```
// Get upcoming free sessions (deans only). --GET Method
http://localhost:8080/session/getupcomingfreesession
```

## Technologies Used

- Go (Golang)
- Fiber - Web framework for Go
- MongoDB - Database for storing user and session data
