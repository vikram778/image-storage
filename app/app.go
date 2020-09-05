package app

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	"github.com/willf/pad"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"image-storage/app/errs"
	"image-storage/db"
	"image-storage/logs"
	"image-storage/migrate"
	"image-storage/migration"
	"image-storage/modules/constant"
	"image-storage/modules/entity"
	"image-storage/out"
	"image-storage/paging"
)

const (
	// ContentType defines Content-Type
	ContentType     = "Content-Type"
	ContentTypeJson = "application/json"
	//MaxPadLeft ...
	MaxPadLeft = 15

	DefaultLimit = "10"

	DocPath = "/api/documentation"

	SQLDatetime = "2006-01-02 15:04:05"
)

type App struct {
	DB         *sqlx.DB
	Logger     *logs.Log
	Mux        *http.ServeMux
	Router     *mux.Router
	HttpClient *http.Client
	Status     int
	RawBody    interface{}
	Port       string
}

func NewApp() App {
	return App{}
}

func test(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("start")
}

func (a *App) Init() {
	var (
		database *db.DB
		err      error
	)

	if database, err = db.New(os.Getenv(constant.EnvDbDriver), os.Getenv(constant.EnvDbOpen)); err != nil {
		return
	}

	a.DB = database.Connection
	a.Router = mux.NewRouter()
	a.Port = os.Getenv("APP_PORT")
	a.Router.HandleFunc("/test", test)
}

func (a *App) Migrate() (err error) {
	conf := &goose.DBConf{
		Env: "default",
		Driver: goose.DBDriver{
			Name:    os.Getenv(constant.EnvDbDriver),
			OpenStr: os.Getenv(constant.EnvDbOpen),
			Dialect: &goose.MySqlDialect{},
		},
	}

	if err = migrate.Process(conf, migration.LocalMigrations); err != nil {
		return err
	}

	return
}

// Listen start listening to the server
func (s *App) Listen() {
	log := logs.New()
	log.Print("Initiating Server")
	log.Print("Server Listening to ", s.Port)
	log.Dump()

	http.ListenAndServe(s.Port, s.Router)
}

func (r *App) FormatException(resource interface{}, err error, errList ...error) {

	var (
		errorString = err.Error()
		values      []interface{}
		mErr        errs.Error
	)

	if len(errList) > 0 {
		for _, item := range errList {
			r.Record("Error", item)
		}
	}
	r.Record("Error", err)

	switch {
	case strings.Contains(errorString, "cannot unmarshal string into Go struct field"):
		err = errors.New(errs.ErrRequestBodyInvalid)
	case strings.Contains(strings.ToLower(errorString), "timeout"):
		err = errors.New(errs.ErrGatewayTimeout)
	}

	mErr, err = errs.GetErrorByCode(errorString)

	if err != nil {
		mErr, _ = errs.GetErrorByCode(errs.ErrCodeNotFound)
	}

	r.Status = mErr.HTTPCode
	r.RawBody = errs.FormateErrorResponse(mErr, values...)

}

func (r *App) Defer(Response http.ResponseWriter) {
	var b bytes.Buffer

	r.Record("Content-Type", Response.Header().Get(ContentType))
	defer func() {
		if r.Status == 0 {
			r.Record("Status", http.StatusInternalServerError)
		} else {
			r.Record("Status", r.Status)
		}

		if r.RawBody != nil {

			if fmt.Sprint(r.RawBody) == "[]" {
				emptyResponse, _ := json.Marshal(make([]int64, 0))
				r.Record("Response", string(emptyResponse))
			} else {
				enc := json.NewEncoder(&b)
				enc.SetEscapeHTML(false)
				enc.Encode(r.RawBody)
				r.Record("Response", strings.Replace(string(b.Bytes()), "\n", "", -1))
			}

		}

		r.Record("End", time.Now().Format(SQLDatetime))
		r.Logger.Dump()

		r.Done(Response)

	}()

	if rec := recover(); rec != nil {
		r.Record("Recovery", fmt.Sprint(rec))
		r.FormatException(r, errors.New(fmt.Sprint(rec)))

		if r.Status == 0 {
			r.FormatException(r, errors.New(errs.ErrInternalAppError))
			out.JSON(Response, http.StatusInternalServerError, r.RawBody)
			return
		}

		if r.RawBody != nil {
			out.JSON(Response, r.Status, r.RawBody)
			return
		}
		out.Status(Response, r.Status)
	}
}

// GetQuery fetches the value from the query string and d if empty
func (r *App) GetQuery(Request *http.Request, key string, d string) string {
	v := Request.URL.Query().Get(key)

	if v == "" {
		return d
	}

	return v
}

// Done will handle the primary response processing
func (r *App) Done(Response http.ResponseWriter) {
	defer func() {
		if recover := recover(); recover != nil {
			r.FormatException(r, errors.New(fmt.Sprint(recover)))
		}
	}()

	body := r.RawBody
	status := r.Status

	if body == nil {
		out.Status(Response, status)
		return
	}

	out.JSON(Response, r.Status, body)
}

func (r *App) GetParams(o interface{}, Response http.ResponseWriter, Request *http.Request) (err error) {
	ct := entity.GetContentType(Request)
	if !entity.ValidContentType(ct) {
		r.Status = http.StatusUnsupportedMediaType
		Response.Header().Set("Accept", ContentTypeJson)
		err = errors.New("Unsupported media type")
		return
	}

	body, _ := ioutil.ReadAll(Request.Body)
	// Restore the io.ReadCloser to its original state
	Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if len(body) < 1 {
		r.Status = http.StatusUnprocessableEntity
		err = errors.New(errs.ErrEmptyBodyContent)
		return
	}

	if entity.CheckJSONCT(ct) {
		err = json.Unmarshal(body, o)
		if err != nil {
			r.Status = http.StatusBadRequest
			err = errors.New(errs.ErrRequestBodyInvalid)
			return
		}
	} else if entity.CheckFormDataCT(ct) {
		var frmInput url.Values
		frmInput, err = entity.ParseForm(ct, Request)
		if err == nil {
			decoder := schema.NewDecoder()
			decoder.SetAliasTag("json")
			err = decoder.Decode(o, frmInput)
			if err != nil {
				r.Status = http.StatusBadRequest
				err = errors.New(errs.ErrRequestBodyInvalid)
				return
			}
		}
	}
	return
}

// Paginate create a pagination
func (r *App) Paginate(Request *http.Request, count int64) {
	var (
		limit, offset int
	)

	limit, _ = strconv.Atoi(r.GetQuery(Request, "limit", DefaultLimit))
	offset, _ = strconv.Atoi(r.GetQuery(Request, "offset", "0"))

	page := paging.NewPaging(r.RawBody, offset, limit, count)
	page.Init(Request)

	r.RawBody = page

}

// Body returns the body from the request
func (r *App) Body(Request *http.Request) []byte {
	body, _ := ioutil.ReadAll(Request.Body)
	// Restore the io.ReadCloser to its original state
	Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

func (r *App) Record(key string, value interface{}) {
	r.Logger.Print(pad.Right(key, MaxPadLeft, " "), value)
}
