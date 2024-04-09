package middlewares

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMyRecoverer(t *testing.T) {

	testTable := []struct {
		name               string
		testHandler        http.HandlerFunc
		expectedStatusCode int
	}{
		{
			name: "without_panic",
			testHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatusCode: 200,
		},
		{
			name: "with_panic",
			testHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			}),
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			w := httptest.NewRecorder()

			MyRecoverer(testCase.testHandler).ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
