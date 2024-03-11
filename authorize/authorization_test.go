package authorize

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chocolate/contrib/authorize/tokenutil"
)

func TestAuthorization(t *testing.T) {
	rep := NewSimpleUserRepository()
	rep.Add(&SimpleUser{
		ID:       "1",
		Username: "test",
		Secret:   "123456",
		Password: "123456",
		Claims:   map[string]string{"foo": "bar"},
	})
	m := tokenutil.NewManager()
	auth := New(rep, m)

	var token string
	{
		form := url.Values{}
		form.Set("username", "test")
		form.Set("password", "ea48576f30be1669971699c09ad05c94")
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("X-Client-Id", "1")
		auth.HTTPHandler().ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fail()
		}
		body := response.Body.String()
		t.Log(response.Code)
		t.Log(body)
		token = body[10 : len(body)-3]
	}

	{
		request := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		if claims, err := auth.ValidateHTTPRequest(request); err != nil {
			t.Error(err)
		} else {
			t.Log(claims.Encode())
		}
	}
}
