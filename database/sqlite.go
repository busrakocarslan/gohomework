package database

import (
	"log"
	"database/sql" // go da sql ile etkileşim kurmak için kullanılan standart kütüphane
	_ "modernc.org/sqlite" // Saf Go SQLite sürücüsü, terminalde go get yap / başında _ olması sadece yan etkileri için kullanılıyor anlamına geliyormuş o olmadan kullanınca hata anlınıyor. 
	//*son iki import birlikte sqlite veritabanına bağlanmak için gerekli olan sürücüyü sağlıyor.
)



func SQLLite() *sql.DB{
	var err error
    // veritabanı aç
	db, err := sql.Open("sqlite", "./test.db")//sqlite ile modern sql olana bağlanıyoruz. test.db db adı
	if err != nil {
		log.Fatal(err)// bağlantı hatası olursa db ile hata mesajı yazılıyor ve programın çalışmasını durduruyor 
	}
    //db.ping dbnin çalışır olup olmadığını kontrol ediyor. ping denen bir sinyali db ye fırlatıyor
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to SQLite database!")
	return db
}