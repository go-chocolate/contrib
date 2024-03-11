package authorize

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/go-chocolate/contrib/authorize/tokenutil"
)

type Request interface {
	GetUsername() string
	GetPassword() string
	GetClientID() string
}

type formRequest struct {
	username string
	password string
	clientID string
}

func (r *formRequest) GetUsername() string {
	return r.username
}
func (r *formRequest) GetPassword() string {
	return r.password
}
func (r *formRequest) GetClientID() string {
	return r.clientID
}

type RequestBuilder func(request *http.Request) (Request, error)

type PasswordVerifier func(pwd, secret, input string) bool

var (
	defaultRequestBuilder RequestBuilder = func(request *http.Request) (Request, error) {
		clientID := request.Header.Get("X-Client-ID")
		switch request.Header.Get("Content-Type") {
		case "application/json":
			m := map[string]any{}
			if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
				return nil, err
			}
			return &formRequest{
				username: m["username"].(string),
				password: m["password"].(string),
				clientID: clientID,
			}, nil

		case "application/x-www-form-urlencoded", "form-data":
			return &formRequest{
				username: request.FormValue("username"),
				password: request.FormValue("password"),
				clientID: clientID,
			}, nil
		default:
			return nil, ErrUnsupportedContentType
		}

	}

	defaultPasswordVerifier PasswordVerifier = func(pwd, secret, input string) bool {
		b := md5.Sum([]byte(pwd + secret))
		expected := hex.EncodeToString(b[:])
		return expected == input
	}
)

type Option func(*Authorization)

func applyOptions(a *Authorization, options []Option) {
	for _, option := range options {
		option(a)
	}
}

func WithRequestBuilder(builder RequestBuilder) Option {
	return func(a *Authorization) {
		a.requestBuilder = builder
	}
}

func WithPasswordVerifier(verifier PasswordVerifier) Option {
	return func(a *Authorization) {
		a.passwordVerifier = verifier
	}
}

type Authorization struct {
	rep              UserRepository
	token            *tokenutil.Manager
	userRepository   UserRepository
	requestBuilder   RequestBuilder
	passwordVerifier PasswordVerifier
}

func New(rep UserRepository, token *tokenutil.Manager, options ...Option) *Authorization {
	a := &Authorization{
		rep:              rep,
		token:            token,
		userRepository:   rep,
		requestBuilder:   defaultRequestBuilder,
		passwordVerifier: defaultPasswordVerifier,
	}
	applyOptions(a, options)
	return a
}

func (a *Authorization) Authorize(ctx context.Context, request Request) (string, error) {
	user, err := a.rep.GetByUsername(ctx, request.GetUsername())
	if err != nil {
		return "", ErrInvalidUsername
	}
	if !a.passwordVerifier(user.GetPassword(), user.GetSecret(), request.GetPassword()) {
		return "", ErrInvalidPassword
	}
	return a.token.GenToken(ctx, user.GetID(), request.GetClientID(), user.GetClaims())
}

func (a *Authorization) AuthorizeFromHTTPRequest(request *http.Request) (string, error) {
	req, err := a.requestBuilder(request)
	if err != nil {
		return "", err
	}
	return a.Authorize(request.Context(), req)
}

func (a *Authorization) HTTPHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		token, err := a.AuthorizeFromHTTPRequest(request)
		if err != nil {
			var message string
			switch err {
			case ErrInvalidUsername, ErrInvalidPassword:
				message = "invalid username or password"
			default:
				message = "system error"
			}
			http.Error(writer, message, http.StatusUnauthorized)
		} else {
			json.NewEncoder(writer).Encode(map[string]interface{}{"token": token})
		}
	}
}

func (a *Authorization) ValidateHTTPRequest(request *http.Request) (tokenutil.Claims, error) {
	return a.token.ValidateHTTPRequest(request)
}

func (a *Authorization) HTTPMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return a.httpMiddleware(next)
	}
}

func (a *Authorization) HTTPMiddlewareFunc() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return a.httpMiddleware(next)
	}
}

func (a *Authorization) httpMiddleware(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if claims, err := a.ValidateHTTPRequest(request); err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
		} else {
			next.ServeHTTP(writer, request.WithContext(claims.WithContext(request.Context())))
		}
	}
}
