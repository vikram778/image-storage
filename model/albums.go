package model

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
	"image-storage/app/errs"
	"time"
)

type Album struct {
	Db         *sqlx.DB    `db:"-" json:"-"`
	ID         null.Int    `db:"id" json:"id"`
	Tittle     null.String `db:"tittle" json:"tittle"`
	ImageCount null.Int    `db:"image_count" json:"image_count"`
	CreatedAt  null.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  null.Time   `db:"updated_at" json:"updated_at"`
}

// NewAlbum function ...
func NewAlbum(db *sqlx.DB) (*Album, error) {
	if db == nil {
		return nil, errors.New("No databse connection")
	}

	return &Album{Db: db}, nil
}

func (me *Album) GetAlbum() (err error) {

	query := `SELECT * FROM albums WHERE tittle = ?`
	err = me.Db.Get(me, query, me.Tittle)
	if err != nil && err != sql.ErrNoRows {
		err = errors.New(errs.ErrInternalDBError)
		return
	}

	return nil
}

func (me *Album) GetAlbumByID() (err error) {

	query := `SELECT * FROM albums WHERE id = ?`
	err = me.Db.Get(me, query, me.ID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.New(errs.ErrInternalDBError)
		return
	}

	return nil
}

func (me *Album) InsertOrUpdate(ok bool) (err error) {
	if ok {
		me.CreatedAt.SetValid(time.Now())
		me.ImageCount.SetValid(0)
		query := `INSERT INTO albums (tittle, image_count, created_at) VALUES (?,?,?)`
		result, errsql := me.Db.Exec(query, me.Tittle, me.ImageCount, me.CreatedAt)
		if errsql != nil {
			err = errors.New(errs.ErrInternalDBError)
			return
		}

		id, _ := result.LastInsertId()
		me.ID.SetValid(id)
	} else {
		me.UpdatedAt.SetValid(time.Now())
		query := `UPDATE albums set  image_count = ? ,updated_at=? WHERE id =?`
		_, errsql := me.Db.Exec(query, me.ImageCount, me.UpdatedAt, me.ID)
		if errsql != nil {
			err = errors.New(errs.ErrInternalDBError)
			return
		}
	}
	return nil
}
