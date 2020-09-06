package app

import (
	"errors"
	"image-storage/app/errs"
	"image-storage/app/resource/api/image"
	"image-storage/filesystem"
	"image-storage/logs"
	"image-storage/model"
	"log"
	"net/http"
	"os"
	"time"
)

func (a *App) PostImage(w http.ResponseWriter, r *http.Request) {
	var (
		albm, _      = model.NewAlbum(a.DB)
		img, _       = model.NewImage(a.DB)
		req          image.PostImageRequest
		res          image.PostImageResponse
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
	//a.Record("Request", strings.Replace(string(a.Body(r)), "\n", "", -1))

	err = a.GetParams(&req, w, r)
	if err != nil {
		a.FormatException(r, err)
		return
	}
	a.Record("Album Name:", req.AlbumTittle)
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err = r.ParseMultipartForm(10 << 20)

	// Multipart form
	form := r.MultipartForm
	if err != nil {
		log.Println(err.Error())
		a.FormatException(r, err)
		return
	}

	files := form.File["image"]
	a.Record("Album Size:", len(files))
	if len(files) > 1 {
		a.FormatException(r, errors.New(errs.ErrMaxLimit))
		return
	}

	// FormFile returns the first file for the given key `image`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, imag, err := r.FormFile("image")
	if err != nil {
		a.Record("Error Retrieving the File:", err.Error())
		a.FormatException(r, err)
		return
	}

	defer file.Close()
	a.Record("Uploaded File:", imag.Filename)
	a.Record("File Size:", imag.Size)
	a.Record("MIME Header:", imag.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	if os.Getenv("ALBUM_FOLDER") != "" {
		album_folder = os.Getenv("ALBUM_FOLDER")
	}

	albumpath := album_folder + "/" + req.AlbumTittle

	filepath := albumpath + "/" + imag.Filename
	err = filesystem.WriteImage(albumpath, filepath, file)
	if err != nil {
		a.Record("Error writing image the File:", err.Error())
		a.FormatException(r, err)
		return
	}

	albm.Tittle.SetValid(req.AlbumTittle)
	err = albm.GetAlbum()
	if err != nil {
		a.FormatException(r, err)
		return
	}

	if !albm.ID.Valid {
		a.FormatException(r, errors.New(errs.ErrAlbumNotExist))
		return
	}

	img.AlbumId.SetValid(albm.ID.Int64)
	img.ImagePath.SetValid(filepath)
	err = img.Insert()
	if err != nil {
		a.FormatException(r, err)
		return
	}

	count := albm.ImageCount.Int64 + 1
	albm.ImageCount.SetValid(count)

	err = albm.InsertOrUpdate(false)
	if err != nil {
		a.FormatException(r, err)
		return
	}

	res.ImageCount = count
	res.AlbumTittle = albm.Tittle.String
	res.ImageName = imag.Filename
	res.ImageID = img.ID.Int64

	a.RawBody = res
	a.Status = http.StatusOK
	return

}
