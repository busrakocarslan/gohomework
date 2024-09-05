package main // Programın ana paketi, kodun hangi paketin içinde olduğunu belirtmek için kullanılıyor.

import ( // ihtiyaç olan kütüphaneleri alıyorum.
	"gohomework/database"
	"gohomework/middleware"
	"gohomework/router"
	"log" // hata ayıklama/loglama
	"os" // Portu çevresel değişkenden almak için

	"github.com/labstack/echo/v4" // echo kullanımı için temel kütüphane
	//ORM-gorm kullanmadım.
	"github.com/swaggo/echo-swagger" // Swagger için gerekli paket
	_ "gohomework/docs" // Swagger dokümantasyon dosyalarını import etmelisin
)

func main() {
	// SQLite veritabanına bağlama
	db := database.SQLLite()
	var err error

	// Var olan tabloyu silme yanlış oluşunca ihtiyacım oldu
	// dropTableQuery := `DROP TABLE IF EXISTS notes;`
	// _, err := db.Exec(dropTableQuery)
	// if err != nil {
		// log.Fatal("Failed to drop table:", err)
	// }

	// Burada tabloyu oluşturuyorum.// if not exists= varsa o tablo atlıyor burayı 
	createTableQuery := `
	  CREATE TABLE IF NOT EXISTS notes (
		   id UUID PRIMARY KEY,
		   title TEXT NOT NULL,
		   content TEXT,		 
		   created_at TEXT NOT NULL,
		   updated_at TEXT NOT NULL
	  )`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Veritabanına başarılı bir şekilde bağlanıldığında
	log.Println("Database connected and table created!")

	rdb := database.UseRedis()

	e := echo.New() // echo framework sınıfından bir instance oluşturuyoruz

	// auth middleware i ekle 
	e.Use(middleware.BasicAuthMidd())

	// Swagger dokümantasyonu endpoint'i
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Router'ı başlat, içerisine parametre olarak Echo'yu ver
	router.NoteRoutes(e, db, rdb)

	// Portu çevresel değişkenden alma, yoksa varsayılan port 9091
	port := os.Getenv("PORT")
	if port == "" {
		port = "9091" // Render'ın portu belirleyememesi durumunda varsayılan bir port ayarlayın
	}

	// Sunucuyu başlat
	e.Logger.Fatal(e.Start(":9091"))
}
