package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  " test #1",
			input: "",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
		{
			name:  " test #2",
			input: "{\"login\": \"user1\",\"password\": \"password1\"}",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Log("", test.name)
		request := httptest.NewRequest(http.MethodPost, "/status", bytes.NewBuffer([]byte(test.input)))
		// создаём новый Recorder
		w := httptest.NewRecorder()
		Register(w, request)

		res := w.Result()
		// проверяем код ответа
		//assert.Equal(t, test.want.code, res.StatusCode)
		// получаем и проверяем тело запроса
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)

		t.Logf(string(resBody))
		t.Log(res.StatusCode)
		require.NoError(t, err)
		assert.Equal(t, test.want.code, res.StatusCode)
		//assert.JSONEq(t, test.want.response, string(resBody))
		//assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
	}
}
