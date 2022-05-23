package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("/users/:id", getUser)
	e.GET("/show", showuser)
	e.POST("/form", formSave)
	e.POST("/formData", formData)
	e.POST("/users", users)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// URL 路徑參數
func getUser(c echo.Context) error {
	// User ID 來自url users:/:id

	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

// 請求參數
func showuser(c echo.Context) error {
	// 從請求參數 user 和 phone 的值

	user := c.QueryParam("user")
	phone := c.QueryParam("phone")
	return c.String(http.StatusOK, "user:"+user+" \nphone:"+phone)
}

// 表單 application/x-www-form-urlencoded
func formSave(c echo.Context) error {

	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+"\nemail:"+email)
}

// 表單 multipart/form-data
func formData(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<p>Thank you! "+name+"</p>")
}

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func users(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, u)
	// c.XML(http.StatusCreated, u)
}
