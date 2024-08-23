package handlers

import ( // virgül yok dikkat et

	"context"
	"database/sql"
	"encoding/json"
	"gohomework/model"
	"log"
	"net/http"
	"time" // cerate-update için

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	// "honnef.co/go/tools/lintcmd/cache"
)

// go daki map= js de obj (key-value)
//? var notes = map[uuid.UUID]model.Note{}// başlangıçta sabit değer ile başlar bir alternatifi make ==> DB KULLANDIKTAN SONRA KALDIRILDI!!!


var validate *validator.Validate

func InitValidation() {
	validate = validator.New()
}

//*********************************************/
// sqlite sonrası yazılan yeni CREATE işlemi 
//?NOT: create func içine func parametresi db *sql.DB, sqlite veri tabanıyla etkileşim kurmamı sağlayan *sql.DB  türünde bir bağlandtı. 
//* echo.handlerFunc echo frameworkunda http handler fonksiyonu.

func CreateNote(db *sql.DB) echo.HandlerFunc {// db ye bağlan sql sorgu çalıştır
	return func(c echo.Context) error {

		note := new(model.Note)// yeni bir note.model instance i 

		// req sonra gelen JSON verisini note a bağlıyorum, bind etmek.
		if err := c.Bind(note); err != nil {//bind edilemezse dönecek hata 
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		// validation konrtolünü burada yapıyorum. Kurallar model de tanımlı
		if err := validate.Struct(note); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		// yeni bir ıd oluşturuluyor uuid.new() ile 
		note.ID = uuid.New()
		// 'CreatedAt' ve 'UpdatedAt' alanlarını, notun oluşturulma zamanına göre ayarlıyor yukarıya tiem i import etti 
		note.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		note.UpdatedAt = note.CreatedAt
		//time.now= go da "an"=> 2024-08-22 14:35:48.123456789 +0300 EEST
		// format standart bir format sağlıyor string dönüyor.

		// SQL sorgu ile veritabanına ekliyorum, query i ayrı tanımlamayıp direk Exec içinde ilk sırada parazntez ile de yazabilirsin. 
		query := `INSERT INTO notes (id, title, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
        _, err := db.Exec(query, note.ID, note.Title, note.Content, note.CreatedAt, note.UpdatedAt) //exec= veriyi yazıyor, create,update,delete işlemlerinde.
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		
		
        // başarılı olursa res json dönüyor
		return c.JSON(http.StatusCreated, note)
	}
}

//**********************************************************/

// sqlite sonrası yazılan yeni read işlemi 
func GetNotes(db *sql.DB, rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {// echo.context req ve res içeriyor
		// Context oluşturuluyor, bu bağlam Redis işlemi için 
		ctx := context.Background()

		// Redis'te cache kontrolü yapıyorum
		cachedNote, err := rdb.Get(ctx,"notes").Result()
		// Eğer Redis'te veri bulunamazsa (redis.Nil hatası), veritabanına sorgu yapıyorum

		if err == redis.Nil {
				// rediste yoksa db den songulanıyor 
				log.Println("DB'den geliyor")// veritabanından alınınca yazacak terminalden kontrol et

				var notes []model.Note // notes adında slice(array) tanımlıyorum, eklenen tüm note lar bunda tutuluyor 
				rows, err := db.Query(`SELECT id, title, content, created_at, updated_at FROM notes`) //db.query ile sql sorgu gönderiyorum ve sonucu rows olarak döndürüyor,hata alırsa err dönecek. içerisine select ile seçtiklerimi getiriyor.
				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}
				defer rows.Close() // get işlemi bittikten sonra rows'un kapatılmasını sağlıyır, DBnin sürekli işlem görmemesi için öenmli!
		
				for rows.Next() {// bu döngü ile sql in her satırı tek tek işleniyor. row.text ile bir sonraki satır varmı bakılıyor.// bu kısım olmasa sadece ilk row işlenir, ve verileri doğru yere atıyor
					var note model.Note
					if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
						return c.JSON(http.StatusInternalServerError, err.Error())
					}// & başlayan yer ile sorguda dönen veriler ilgili yere aktarılıyor. 
					notes = append(notes, note)// note, notes' a ekleniyor. 
				}
				  // Notes JSON formatında seri hale getiriliyor
		          notesJSON,_ := json.Marshal(notes)
				  // Veri Redis'te cache olarak saklanıyor, set ediliyor (5 dakika)
	            	rdb.Set(ctx,"notes",notesJSON, 5*time.Minute)

					return c.JSON(http.StatusOK, notes)

		}else if err !=nil {
			// Redis'ten veri alırken başka bir hata oluşursa
			return c.JSON(http.StatusInternalServerError,err.Error())

		}
		log.Println("Redisten geliyor")// redisten alınırsa yazacak
		var notes []model.Note
			// Redis'ten veri başarıyla alındıysa, JSON formatından Go struct'ına dönüştürmeliyim
			json.Unmarshal([]byte(cachedNote),&notes)
			// Cache'ten alınan not JSON formatında geri döndürülüyor
			return c.JSON(http.StatusOK, notes)
	}
}

//*****************************************************************/
// tek okuma işlemi 
func ReadNote(db *sql.DB, rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Context oluşturuluyor, bu bağlam Redis işlemi için 
		ctx := context.Background()
		id, err:= uuid.Parse(c.Param("id"))// c.param ile url deki id alınıyor, parse işlemi ile alınann string olduğundan uuid ye dönüştürülüyor. 
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid UUID")// uuid karakter sayısı yanlış ise giriyor
		}
		
		// uuid kontrolü sonrası Redis'te cache kontrolü yapıyorum
		cachedNote, err := rdb.Get(ctx, id.String()).Result()
		

		// Eğer Redis'te veri bulunamazsa (redis.Nil hatası), veritabanına sorgu yapıyorum

		if err == redis.Nil {
			// rediste yoksa db den songulanıyor daha öncekki kod direk bu func içine girdi 
			log.Println("DB'den geliyor")// veritabanından alınınca yazacak terminalden kontrol et

			var note model.Note // model note dan id ile sorguladığım note da tutuluyor

			query := `SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?`// /? yerine sorgu başladığında  id geliyor
			// QueryRow sönen ilk row u alıyor
			if err := db.QueryRow(query, id).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
				if err == sql.ErrNoRows {
					return c.JSON(http.StatusNotFound, "Note not found")
				}
				return c.JSON(http.StatusInternalServerError, err.Error())

		}
        // Not JSON formatında seri hale getiriliyor
		noteJSON,_ := json.Marshal(note)
		// Veri Redis'te cache olarak saklanıyor, set ediliyor (5 dakika)
		rdb.Set(ctx,id.String(),noteJSON, 5*time.Minute)

		// Veritabanından alınan not, JSON formatında geri dönsün
		return c.JSON(http.StatusOK, note)	
		

	} else if err !=nil {
		// Redis'ten veri alırken başka bir hata oluşursa
			return c.JSON(http.StatusInternalServerError,err.Error())
		}
		log.Println("Redisten geliyor")// redisten alınıca yazacak 

		var note model.Note // model note dan id ile sorguladığım note da tutuluyor
		

		// Redis'ten veri başarıyla alındıysa, JSON formatından Go struct'ına dönüştürmeliyim
		json.Unmarshal([]byte(cachedNote),&note)
		// Cache'ten alınan not JSON formatında geri döndürülüyor
		return c.JSON(http.StatusOK, note)
		

		// scan= veriyi okuyor, note a atanıyor, get işleminde 

	}
}

//*****************************************************************/
// update işlemi 
func UpdateNote(db *sql.DB, rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid UUID")
		}
		note := new(model.Note)
		// veri alınıp note nesnesine bind ediliyor 
		if err := c.Bind(note); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		// validation işlemi yapılıyor. 

		if err := validate.Struct(note); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		note.ID=id // urldeki idyi note ID ye atanıp güncelenecek nesne belirleme 
		note.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		//time.now= go da "an"=> 2024-08-22 14:35:48.123456789 +0300 EEST
		// format standart bir format sağlıyor string dönüyor.


		query := `UPDATE notes SET title = ?, content = ?, updated_at = ? WHERE id = ?`// sql sorgusu ile idye göre güncellenecek alanlar
        _, err = db.Exec(query, note.Title, note.Content, note.UpdatedAt, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		// Context oluşturuluyor, bu bağlam Redis işlemi için //! burada dikkat edilecek husus tüm update işlemini yukarıda yaptıktan sonra son return dan önce redis işlemini ekliyoruz. 
		ctx := context.Background()
	    // Not JSON formatında seri hale getiriliyor
		noteJSON,_ := json.Marshal(note)
		// Veri Redis'te cache olarak saklanıyor, set ediliyor (5 dakika)
		rdb.Set(ctx,id.String(),noteJSON, 5*time.Minute)

		return c.JSON(http.StatusOK, note)
	}
}
//*****************************************************************/
// delete  işlemi 

func DeleteNote(db *sql.DB,rdb *redis.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid UUID")
		}

		query := `DELETE FROM notes WHERE id = ?` // sql saldırı riskini ortadan kaldırmak adına ?(placeholder deniyormmuş) kullanılıyor. c.param(id) ile concat yapılmıyor 
		_, err = db.Exec(query, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
        // Context oluşturuluyor, bu bağlam Redis işlemi için
		ctx := context.Background()
		// dikkate et set değil del işlemi
		rdb.Del(ctx, id.String())

		return c.NoContent(http.StatusNoContent)
	}
}

// ilk db siz kodlar 
// CreateNote: Yeni bir not oluşturur
// func CreateNote(c echo.Context) error{
// 	note :=new(model.Note)// modelden oluşturulan bir note objesi ataması yapılıyor-- api den gelen verileri saklayacak 
// 	if err := c.Bind(note); err !=nil{// c.Bind(node)= HTTp isteğindeki veriyi node nesnesine bağlıyor, hata gelirse err de saklanacak 
// 		return c.JSON(http.StatusBadRequest,err.Error())
// 	}// hataya girerse içene gireceği kod

// 	note.ID=uuid.New()// yeni bir ıd oluşturuluyor uuid.new() ile 
// 	notes[note.ID]= *note// /*note=note burada notes mapine ekleniyor yani note:ID= key *note=value bu sayede UUID ile erişebiliyorum bu note a 
	
// 	return c.JSON(http.StatusCreated,note)// create işlemi başarılı sonuçlanınca yeni oluşturulan note u json formatında dönüyor
// }

// GetNotes: Tüm notları getir 
// func GetNotes(c echo.Context) error {
// 	if len(notes) == 0 {// notes map i nin uzunluğu 0 değilse girmeyecek
// 		return c.JSON(http.StatusNotFound, "No notes found")
// 	}
// 	return c.JSON(http.StatusOK, notes)
// }