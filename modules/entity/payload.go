package entity

import (
	"bytes"
	"github.com/frozentech/go-tools/array"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// payload const declaration
const (
	CT          = "Content-Type"
	CTJson      = "application/json"
	CTPlain     = "text/plain"
	CTFormData  = "multipart/form-data"
	CTUrlEncode = "application/x-www-form-urlencoded"
)

// Helper variable declaration
var (
	JSONSupportCT = []string{CTJson, CTPlain}
	FormSupportCT = []string{CTFormData, CTUrlEncode}
)

// GetContentType ...
func GetContentType(req *http.Request) (ct string) {
	ct = req.Header.Get(CT)
	if flg := strings.Contains(ct, ";"); flg {
		ctWithBountry := strings.Split(ct, ";")
		ct = ctWithBountry[0]
	}
	return
}

// ValidContentType ...
func ValidContentType(ct string) bool {
	if ct == "" {
		return true
	} else if CheckJSONCT(ct) {
		return true
	} else if CheckFormDataCT(ct) {
		return true
	}
	return false
}

// CheckJSONCT - check json related content type
func CheckJSONCT(ct string) bool {
	exist, _ := array.InArray(ct, JSONSupportCT)
	return exist
}

// CheckFormDataCT - check form data related content type
func CheckFormDataCT(ct string) bool {
	exist, _ := array.InArray(ct, FormSupportCT)
	return exist
}

// ParseForm ...
func ParseForm(ct string, r *http.Request) (res url.Values, err error) {

	if r.Body != nil {
		// read all bytes from content body and create new stream using it.
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// create new request for parsing the body
		req2, _ := http.NewRequest(r.Method, r.URL.String(), bytes.NewReader(bodyBytes))
		req2.Header = r.Header

		if ct == CTUrlEncode {
			if err = req2.ParseForm(); err != nil {
				return
			}
			res = req2.PostForm
		} else if ct == CTFormData {
			if err = req2.ParseMultipartForm(200000); err != nil {
				return
			}
			res = req2.PostForm
		}

	}

	return
}
