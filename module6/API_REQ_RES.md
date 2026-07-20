# Student API - Request & Response Documentation

## Overview

This API provides CRUD operations for managing students with JWT-based authentication.

### Base URL

```
http://localhost:8080
```

### Authentication

Protected endpoints require:

```
Authorization: Bearer <JWT_TOKEN>
```

JWT tokens are generated using the `/login` endpoint.

---

# Authentication API

## 1. Login

Generate JWT token.

### Endpoint

```
POST /login
```

### Authentication

No authentication required.

---

## Request

### Headers

```
Content-Type: application/json
```

### Body

```json
{
    "username": "admin",
    "password": "admin123"
}
```

---

## Success Response

### Status

```
200 OK
```

### Response Body

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## Error Response

### Invalid Credentials

Status:

```
401 Unauthorized
```

Response:

```json
{
    "message": "Invalid Credentials"
}
```

---

# Student APIs

All student endpoints require JWT authentication.

---

# 2. Get All Students

Returns all students.

## Endpoint

```
GET /students
```

---

## Headers

```
Authorization: Bearer <JWT_TOKEN>
```

---

## Success Response

Status:

```
200 OK
```

Response:

```json
[
    {
        "id": 1,
        "name": "Alice",
        "age": 20,
        "email": "alice@example.com"
    },
    {
        "id": 2,
        "name": "Bob",
        "age": 21,
        "email": "bob@example.com"
    }
]
```

---

# 3. Get Student By ID

Returns a single student.

## Endpoint

```
GET /students/{id}
```

Example:

```
GET /students/1
```

---

## Headers

```
Authorization: Bearer <JWT_TOKEN>
```

---

## Success Response

Status:

```
200 OK
```

Response:

```json
{
    "id": 1,
    "name": "Alice",
    "age": 20,
    "email": "alice@example.com"
}
```

---

## Error Response

Status:

```
404 Not Found
```

Response:

```json
{
    "message": "Student not found"
}
```

---

# 4. Create Student

Creates a new student.

## Endpoint

```
POST /students
```

---

## Headers

```
Content-Type: application/json

Authorization: Bearer <JWT_TOKEN>
```

---

## Request Body

```json
{
    "name": "Charlie",
    "age": 22,
    "email": "charlie@example.com"
}
```

---

## Success Response

Status:

```
200 OK
```

Response:

```json
{
    "id": 3,
    "name": "Charlie",
    "age": 22,
    "email": "charlie@example.com"
}
```

---

## Validation Error

Status:

```
400 Bad Request
```

Response:

```json
{
    "message": "Name and Email required"
}
```

---

# 5. Update Student

Updates an existing student.

## Endpoint

```
PUT /students/{id}
```

Example:

```
PUT /students/1
```

---

## Headers

```
Content-Type: application/json

Authorization: Bearer <JWT_TOKEN>
```

---

## Request Body

```json
{
    "name": "Alice Updated",
    "age": 25,
    "email": "alice.updated@example.com"
}
```

---

## Success Response

Status:

```
200 OK
```

Response:

```json
{
    "id": 1,
    "name": "Alice Updated",
    "age": 25,
    "email": "alice.updated@example.com"
}
```

---

## Error Response

Status:

```
404 Not Found
```

Response:

```json
{
    "message": "Student not found"
}
```

---

# 6. Delete Student

Deletes a student.

## Endpoint

```
DELETE /students/{id}
```

Example:

```
DELETE /students/1
```

---

## Headers

```
Authorization: Bearer <JWT_TOKEN>
```

---

## Success Response

Status:

```
200 OK
```

Response:

```json
{
    "message": "Student deleted"
}
```

---

# Authentication Flow

```
Client
  |
  |
  | POST /login
  |
  | username + password
  |
  v
Server
  |
  |
  | Validate User
  |
  | Generate JWT
  |
  v
Client receives token


Client
  |
  |
  | Authorization: Bearer JWT
  |
  v
Protected API
  |
  |
  | Validate JWT
  |
  v
Student Data
```

---

# Postman Collection Examples

## Login Request

```
POST http://localhost:8080/login
```

Body:

```json
{
    "username":"admin",
    "password":"admin123"
}
```

---

## Get Students Request

```
GET http://localhost:8080/students
```

Headers:

```
Authorization: Bearer eyJhbGc...
```

---

## Create Student Request

```
POST http://localhost:8080/students
```

Headers:

```
Authorization: Bearer eyJhbGc...
Content-Type: application/json
```

Body:

```json
{
    "name":"David",
    "age":23,
    "email":"david@example.com"
}
```

---

# HTTP Status Codes

| Status Code | Meaning                |
| ----------- | ---------------------- |
| 200         | Request successful     |
| 400         | Invalid request        |
| 401         | Missing or invalid JWT |
| 404         | Resource not found     |
| 500         | Server error           |

---

# JWT Token Structure

A JWT contains three parts:

```
HEADER.PAYLOAD.SIGNATURE
```

Example:

```
xxxxx.yyyyy.zzzzz
```

---

## Header

```json
{
    "alg":"HS256",
    "typ":"JWT"
}
```

---

## Payload

```json
{
    "username":"admin",
    "exp":1720000000
}
```

---

## Signature

Generated using:

```
JWT Secret Key
```

---

# Current API Features

✅ JWT Authentication
✅ Bearer Token Support
✅ Login Endpoint
✅ Student CRUD
✅ Middleware
✅ JSON Request/Response
✅ Token Expiration
✅ Input Validation
✅ HTTP Status Handling

---

# Future Production Improvements

* Database integration
* User registration
* bcrypt password hashing
* Refresh tokens
* Role-based authorization
* Environment variable secrets
* PostgreSQL/MySQL storage
* Swagger/OpenAPI documentation
* Docker deployment
