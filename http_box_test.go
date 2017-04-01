package packr

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HTTPBox(t *testing.T) {
	r := require.New(t)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(testBox))

	req, err := http.NewRequest("GET", "/hello.txt", nil)
	r.NoError(err)

	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	r.Equal(200, res.Code)
	testFile, err := ioutil.ReadFile(filepath.Join("fixtures", "hello.txt"))
	if err != nil {
		t.Error("Cannot read test file: ", err)
	}
	r.Equal(string(testFile), res.Body.String())
}
