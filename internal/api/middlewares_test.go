package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodMiddleware(t *testing.T) {
	cases := []struct {
		AllowMethod string
		CallMethod  string
		IsValid     bool
	}{
		{AllowMethod: "GET", CallMethod: "GET", IsValid: true},
		{AllowMethod: "GET", CallMethod: "POST", IsValid: false},
		{AllowMethod: "POST", CallMethod: "PUT", IsValid: false},
	}

	next := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.CallMethod, "/", nil)
		w := httptest.NewRecorder()

		fn := methodMiddleware(tc.AllowMethod, next)
		fn(w, req)
		resp := w.Result()

		data, _ := ioutil.ReadAll(resp.Body)

		if tc.IsValid {
			assert.Equal(t, resp.StatusCode, http.StatusOK)
			assert.Equal(t, string(data), "OK")
		} else {
			assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)
			assert.Equal(t, string(data), http.StatusText(http.StatusMethodNotAllowed)+"\n")
		}
	}

}
