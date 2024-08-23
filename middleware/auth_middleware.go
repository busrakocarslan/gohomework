package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BasicAuthMidd() echo.MiddlewareFunc{
	return middleware.BasicAuth(func (username, password string,c echo.Context)(bool,error)  {
		if username == "admin" && password =="password" {
			return true,nil
		}
		return false,nil
		
	})
}

/*
1-go get github.com/labstack/echo/v4/middleware@v4.12.0

2- clean kod adına main.go da değil ayrı bir middleware dosyasında yazdım router da tüm CRUD işlemlerine ayrı ayrı eklemek yerine main de router üzerinde middleware olduğundan use ile kullandım. 

3- 



*/