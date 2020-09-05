package out

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// CTypeJSON defines Content-Type for application/json
	CTypeJSON = "application/json"
	// CTypeText defines Content-Type for text/plain
	CTypeText = "text/plain"
)

// Text sends a Text response body
func Text(r http.ResponseWriter, code int, content string) {
	Output(r, code, CTypeText, []byte(content))
}

// JSON sends a JSON response body
func JSON(r http.ResponseWriter, code int, content interface{}) {
	if fmt.Sprint(content) == "[]" {
		emptyResponse, _ := json.Marshal(make([]int64, 0))
		Output(r, code, CTypeJSON, emptyResponse)
		return
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.Encode(content)
	Output(r, code, CTypeJSON, b.Bytes())
}

// Status sends an empty response but requires a proper HTTP Status Code
func Status(r http.ResponseWriter, code int) {
	Output(r, code, CTypeJSON, nil)
}

// Output sets a full HTTP output detail
func Output(r http.ResponseWriter, code int, ctype string, content []byte) {
	r.Header().Set("Content-Type", ctype)
	r.WriteHeader(code)
	r.Write(content)
}
