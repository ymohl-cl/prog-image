package image

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Image :_
type Image struct {
	session *mgo.Session
}

// New provide a new Image with driver MongoDB
func New(s *mgo.Session) *Image {
	return &Image{session: s}
}

// UploadImage save image and return id attached
func (i Image) UploadImage(c echo.Context) error {
	// Get image
	image, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	// source
	src, err := image.Open()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer src.Close()

	// create image on bdd
	dst, err := i.session.DB("test").GridFS("fs").Create(image.Filename)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	// save content on bdd
	if _, err = io.Copy(dst, src); err != nil {
		c.Logger().Error(err)
		return err
	}

	// get id image
	id := dst.Id()

	if err = dst.Close(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusExpectationFailed, id)
	}

	// return id
	return c.JSON(http.StatusOK, id)
}

// GetImage return image to download it by id
func (i Image) GetImage(c echo.Context) error {
	var file *os.File
	// Get ID
	id := c.FormValue("id")

	// Get image by id
	src, err := i.session.DB("test").GridFS("fs").OpenId(bson.ObjectIdHex(id))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	// Create file to send it
	if file, err = os.Create(id); err != nil {
		c.Logger().Error(err)
		return err
	}
	defer file.Close()
	// delete file
	defer deleteFile(id)

	// Write on the new file
	if _, err = io.Copy(file, src); err != nil {
		c.Logger().Error(err)
		return err
	}

	// send it
	return c.Attachment(id, src.Name())
}

func deleteFile(path string) {
	if err := os.Remove(path); err != nil {
		panic(err)
	}
	return
}
