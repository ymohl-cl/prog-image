package main

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/labstack/echo"
	"github.com/ymohl-cl/prog-image/image"
)

func defineRouting(e *echo.Echo, i *image.Image) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, WORLD by get")
	})
	e.GET("/image", i.GetImage)
	e.POST("/image", i.UploadImage)
}

func main() {
	var e *echo.Echo
	var session *mgo.Session
	var err error

	// connect to bdd MongoDB with GridFS support
	session, err = mgo.Dial("mongo")
	if err != nil {
		panic(err)
	}

	// get image stucture
	image := image.New(session)
	// initialize routing api
	e = echo.New()
	defineRouting(e, image)

	// start the server
	e.Logger.Fatal(e.Start(":8000"))
}
