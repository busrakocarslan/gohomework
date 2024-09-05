## https://gohomework.onrender.com/swagger/index.html


## Notes API - Go & Echo Project

This is a simple Notes API project built using the Go and Echo framework. Basic Authentication has been implemented in the project. Once logged in with admin credentials, CRUD (Create, Read, Update, Delete) operations can be performed on the notes. Each note has a unique UUID and stores created_at and updated_at timestamps. The data is stored in SQLite, and Redis is used for performance improvements.

## Features
Authentication: Log in with "admin" and "password".<br>
Unique UUID: Each note is assigned a unique UUID.<br>
CRUD Operations: Perform Create, Read, Update, and Delete operations on notes.<br>
Timestamps: Each note records created_at and updated_at timestamps.<br>
Redis Integration: Redis is used to improve performance and reduce database load.<br>
SQLite: Data is stored in a lightweight SQLite database.<br>
No ORM: Direct SQL queries are used for database operations without any ORM.<br>

## Technologies Used
Go<br>
Echo<br>
SQLLiTE<br>
Redis <br>
UUID <br>
Basic Authentication ("admin"-"password")<br>