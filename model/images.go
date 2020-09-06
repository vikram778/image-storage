package model

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
	"image-storage/app/errs"
	"time"
)

type Image struct {
	Db        *sqlx.DB    `db:"-" json:"-"`
	ID        null.Int    `db:"id" json:"id"`
	AlbumId   null.Int    `db:"album_id" json:"album_id"`
	ImagePath null.String `db:"image_path" json:"image_path"`
	CreatedAt null.Time   `db:"created_at" json:"created_at"`
}

// NewImage function ...
func NewImage(db *sqlx.DB) (*Image, error) {
	if db == nil {
		return nil, errors.New("No databse connection")
	}

	return &Image{Db: db}, nil
}

func (me *Image) GetImage() (err error) {

	query := `SELECT * FROM images WHERE id= ?`
	err = me.Db.Get(me, query, me.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		err = errors.New(errs.ErrInternalDBError)
		return
	}

	return nil
}

func (me *Image) Insert() (err error) {
	me.CreatedAt.SetValid(time.Now())
	query := `INSERT INTO images (album_id, image_path, created_at) VALUES (?,?,?)`
	result, errsql := me.Db.Exec(query, me.AlbumId, me.ImagePath, me.CreatedAt)
	if errsql != nil {
		err = errors.New(errs.ErrInternalDBError)
		return
	}

	id, _ := result.LastInsertId()
	me.ID.SetValid(id)

	return nil
}

func (me *Image) DeleteImage() (err error) {

	query := `DELETE FROM images WHERE id = ?`
	_, err = me.Db.Exec(query, me.ID)
	if err != nil {
		err = errors.New(errs.ErrInternalDBError)
		return
	}

	return nil
}
