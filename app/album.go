package app

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"image-storage/app/errs"
	"image-storage/app/resource/api/album"
	"image-storage/filesystem"
	"image-storage/kafka/producer"
	"image-storage/logs"
	"image-storage/model"
	"net/http"
	"strings"
	"time"
)

func (a *App) PostAlbum(w http.ResponseWriter, r *http.Request) {

	var (
		albm, _           = model.NewAlbum(a.DB)
		req               album.PostAlbumRequest
		res               album.PostAlbumResponse
		kafkanotification map[string]interface{}
		err               error
		album_folder      string
	)

	defer func() {
		a.Defer(w)
	}()

	a.Logger = logs.New()
	a.Record("Start", time.Now().Format(SQLDatetime))

	a.Record("Resource", "album")
	a.Record("Method", r.Method)
	a.Record("URL", r.URL.String())
	a.Record("Request", strings.Replace(string(a.Body(r)), "\n", "", -1))

	err = a.GetParams(&req, w, r)
	if err != nil {
		a.FormatException(r, err)
		return
	}

	album_folder = a.GetAlbumsDir()

	albumpath := album_folder + "/" + req.AlbumTittle
	if exist, _ := filesystem.Exist(albumpath); exist {
		a.Record("Error Album Doesnt exist:", req.AlbumTittle)
		a.FormatException(r, errors.New(errs.ErrAlbumExist))
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

	kbyt, _ := json.Marshal(res)
	json.Unmarshal(kbyt, &kafkanotification)
	kafkanotification["topic"] = PostAlbumTopic

	producer.Jobs <- kafkanotification

	a.RawBody = res
	a.Status = http.StatusOK

	return
}

func (a *App) DeleteAlbum(w http.ResponseWriter, r *http.Request) {

	var (
		albm, _           = model.NewAlbum(a.DB)
		img, _            = model.NewImage(a.DB)
		res               album.DeleteAlbumResponse
		kafkanotification map[string]interface{}
		err               error
		album_folder      string
	)

	defer func() {
		a.Defer(w)
	}()

	a.Logger = logs.New()
	a.Record("Start", time.Now().Format(SQLDatetime))

	a.Record("Resource", "album")
	a.Record("Method", r.Method)
	a.Record("URL", r.URL.String())

	vars := mux.Vars(r)
	tittle, ok := vars["tittle"]
	if !ok {
		a.FormatException(r, errors.New("invalid tittle"))
		return
	}

	albm.Tittle.SetValid(tittle)
	err = albm.GetAlbum()
	if err != nil {
		a.FormatException(r, err)
		return
	}

	if !albm.ID.Valid {
		a.FormatException(r, errors.New(errs.ErrAlbumNotExist))
		return
	}

	album_folder = a.GetAlbumsDir()

	albumpath := album_folder + "/" + tittle
	exist, _ := filesystem.Exist(albumpath)
	if exist {
		err = filesystem.DeleteDir(albumpath)
		if err != nil {
			a.FormatException(r, err)
		}
	}

	img.AlbumId.SetValid(albm.ID.Int64)
	err = img.DeleteImagesByAlbumID()
	if err != nil {
		a.FormatException(r, err)
	}

	err = albm.DeleteAlbum()
	if err != nil {
		a.FormatException(r, err)
	}

	res.AlbumTittle = albm.Tittle.String
	res.Message = "album deleted successfully"

	kbyt, _ := json.Marshal(res)
	json.Unmarshal(kbyt, &kafkanotification)
	kafkanotification["topic"] = DeleteAlbumTopic

	producer.Jobs <- kafkanotification

	a.RawBody = res
	a.Status = http.StatusOK

	return

}
