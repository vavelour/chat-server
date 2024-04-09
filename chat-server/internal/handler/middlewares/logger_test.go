package middlewares

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMyLogger(t *testing.T) {

	testTable := []struct {
		name               string
		infoMessage        string
		method             string
		expectedStatusCode string
		url                string
		message            string
		logOutput          bytes.Buffer
		testHandler        http.HandlerFunc
	}{
		{
			name:               "ok",
			infoMessage:        "Request processed successfully.",
			method:             "METHOD=GET",
			expectedStatusCode: "STATUS=200",
			url:                "URL=/",
			message:            "MESSAGE=OK",
			testHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
		{
			name:               "bad_request",
			infoMessage:        "Bad Request.",
			method:             "METHOD=GET",
			expectedStatusCode: "STATUS=400",
			url:                "URL=/",
			message:            "MESSAGE=ERROR",
			testHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("ERROR"))
			}),
		},
		{
			name:               "server_error",
			infoMessage:        "Server error occurred.",
			method:             "METHOD=GET",
			expectedStatusCode: "STATUS=500",
			url:                "URL=/",
			message:            "MESSAGE=SERVER_OCCURRED",
			testHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("SERVER_OCCURRED"))
			}),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			logrus.SetOutput(&testCase.logOutput)

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			w := httptest.NewRecorder()

			MyLogger(testCase.testHandler).ServeHTTP(w, req)

			assert.Contains(t, testCase.logOutput.String(), testCase.infoMessage)
			assert.Contains(t, testCase.logOutput.String(), testCase.method)
			assert.Contains(t, testCase.logOutput.String(), testCase.expectedStatusCode)
			assert.Contains(t, testCase.logOutput.String(), testCase.url)
			assert.Contains(t, testCase.logOutput.String(), testCase.message)

			testCase.logOutput.Reset()
		})
	}
}
