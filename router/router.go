package router

import (
	"database/sql" // SQL veritabanı işlemleri için gerekli kütüphane

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	// "gorm.io/gorm"
	"gohomework/handlers"
)

// echo, sql ve redis'e birlikte bağlanıyor
func NoteRoutes(e *echo.Echo, db *sql.DB, rdb *redis.Client){
	handlers.InitValidation() // validation ı başlatıyor. 
    // veritanabı işlemlerinin gerçekleşebilmesi için sonuna db eklenmeli.
	e.POST("/notes",handlers.CreateNote(db))
	e.GET("/notes",handlers.GetNotes(db, rdb))
	e.GET("/notes/:id",handlers.ReadNote(db,rdb))
	e.PUT("/notes/:id",handlers.UpdateNote(db,rdb))
	e.DELETE("/notes/:id",handlers.DeleteNote(db,rdb))
}