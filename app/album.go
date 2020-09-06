package app

import (
	"errors"
	"image-storage/app/errs"
	"image-storage/app/resource/api/album"
	"image-storage/filesystem"
	"image-storage/logs"
	"image-storage/model"
	"net/http"
	"os"
	"strings"
	"time"
)

func (a *App) PostAlbum(w http.ResponseWriter, r *http.Request) {

	var (
		albm, _      = model.NewAlbum(a.DB)
		req          album.PostAlbumRequest
		res          album.PostAlbumResponse
		err          error
		album_folder string
	)

	defer func() {
		a.Defer(w)
	}()

	a.Logger = logs.New()
	a.Record("Start", time.Now().Format(SQLDatetime))

	a.Record("Resource", "story")
	a.Record("Method", r.Method)
	a.Record("URL", r.URL.String())
	a.Record("Request", strings.Replace(string(a.Body(r)), "\n", "", -1))

	err = a.GetParams(&req, w, r)
	if err != nil {
		a.FormatException(r, err)
		return
	}

	if os.Getenv("ALBUM_FOLDER") != "" {
		album_folder = os.Getenv("ALBUM_FOLDER")
	}

	albumpath := album_folder + "/" + req.AlbumTittle
	if exist, _ := filesystem.Exist(albumpath); exist {
		a.Record("Error Album Doesnt exist:", req.AlbumTittle)
		a.FormatException(r, errors.New(errs.ErrAlbumNotExist))
		return
	}

	if err = filesystem.Mkdir(albumpath); err != nil {
		a.FormatException(r, err)
		return
	}

	albm.Tittle.SetValid(req.AlbumTittle)
	err = albm.InsertOrUpdate(true)
	if err != nil {
		a.FormatException(r, err)
		return
	}

	res.AlbumTittle = req.AlbumTittle
	res.Message = "Album created sucessfully"

	a.RawBody = res
	a.Status = http.StatusOK

	return
}
