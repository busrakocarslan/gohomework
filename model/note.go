package model // model yapılarının bulunduğu paket oluyor, bu kısımda istediğimiz adı verebiliyoruz. 

 

//  benzersiz uuid kullanmak için önce terminalde go get yapıp import et go.mod dosyasına gelmişmi kontrol et
 import "github.com/google/uuid"

//note modelini oluştur
//alpha= sadece harf içermeli

type Note struct{
	 ID      uuid.UUID `json:"id"`	 
	 Title   string `json:"title" validate:"required,min=3"`
	 Content string `json:"content"`
	 CreatedAt string    `json:"created_at"`
	 UpdatedAt string    `json:"updated_at"`
	
}