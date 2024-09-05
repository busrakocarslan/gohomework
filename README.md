## Notes API - Go & Echo Project

This is a simple Notes API project built using the Go and Echo framework. Basic Authentication has been implemented in the project. Once logged in with admin credentials, CRUD (Create, Read, Update, Delete) operations can be performed on the notes. Each note has a unique UUID and stores created_at and updated_at timestamps. The data is stored in SQLite, and Redis is used for performance improvements.

## Features
Authentication: Log in with "admin" and "password".
Unique UUID: Each note is assigned a unique UUID.
CRUD Operations: Perform Create, Read, Update, and Delete operations on notes.
Timestamps: Each note records created_at and updated_at timestamps.
Redis Integration: Redis is used to improve performance and reduce database load.
SQLite: Data is stored in a lightweight SQLite database.
No ORM: Direct SQL queries are used for database operations without any ORM.

## Technologies Used
Go
Echo
SQLLiTE
Redis 
UUID 
Basic Authentication ("admin"-"password")