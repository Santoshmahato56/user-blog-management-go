# User and Blog Management System

A Go web application that implements user management with dynamic roles and permissions, as well as blog creation and posting functionality.

## Features

- User authentication (register, login, logout)
- Role-based access control (RBAC)
- Dynamic permission management
- Blog creation and management
- Permission-based authorization for blog operations

## Tech Stack

- Go 1.21+
- Gin Web Framework
- GORM (Object Relational Mapper)
- JWT for authentication
- SQLite/PostgreSQL/MySQL database options


## Setup and Installation

1. Clone the repository
2. Create a `.env` file (see `.env.example`)
3. Install dependencies:
   ```
   go mod tidy
   ```
4. Run the application:
   ```
   go run main.go
   ```

## API Endpoints

### Authentication

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Log in
- `GET /auth/me` - Get current user info

### Blogs

- `GET /blogs` - List all published blogs
- `GET /blogs/:id` - Get a blog by ID
- `GET /blogs/user/:user_id` - List blogs by user
- `POST /blogs` - Create a new blog (requires authentication)
- `PUT /blogs/:id` - Update a blog (requires authentication)
- `DELETE /blogs/:id` - Delete a blog (requires authentication)

## Role-Based Permissions

The system has two default roles:

1. **Admin** - Has all permissions
2. **User** - Has read-only permissions

Permissions include:
- `create_blog` - Can create blog posts
- `read_blog` - Can read blog posts
- `update_blog` - Can update blog posts
- `delete_blog` - Can delete blog posts
- `create_user` - Can create users
- `read_user` - Can read user information
- `update_user` - Can update user information
- `delete_user` - Can delete users
