package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pedrocelso/go-rest-service/lib/http/controllers"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestBasicAuthSucceed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.Use(controllers.CheckJWT("exevo gran mas vis"))

	r.GET("/login", func(c *gin.Context) {
		c.String(200, "authorized")
	})

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiUGVkcm8gQ29zdGEiLCJlbWFpbCI6InBlZHJvY2Vsc29udW5lc0BnbWFpbC5jb20ifQ.Xd7vOBrbMuoxEwZpxWDH4bVVvPIFj60fzlt8Kd9Krwk")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "authorized", w.Body.String())
}

func TestWrongSecret(t *testing.T) {
	req, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.Use(controllers.CheckJWT("gerere"))

	r.GET("/login", func(c *gin.Context) {
		c.String(200, "authorized")
	})

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiUGVkcm8gQ29zdGEiLCJlbWFpbCI6InBlZHJvY2Vsc29udW5lc0BnbWFpbC5jb20ifQ.Xd7vOBrbMuoxEwZpxWDH4bVVvPIFj60fzlt8Kd9Krwk")
	r.ServeHTTP(w, req)

	assert.Equal(t, "signature is invalid\n", w.Body.String())
}

func TestNoAuthorization(t *testing.T) {
	req, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.Use(controllers.CheckJWT("gerere"))

	r.GET("/login", func(c *gin.Context) {
		c.String(200, "authorized")
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, "Required authorization token not found\n", w.Body.String())
}

func TestBadAuthorization(t *testing.T) {
	req, _ := http.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.Use(controllers.CheckJWT("gerere"))

	r.GET("/login", func(c *gin.Context) {
		c.String(200, "authorized")
	})

	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiUGVkcm8gQ29zdGEiLCJlbWFpbCI6InBlZHJvY2Vsc29udW5lc0BnbWFpbC5jb20ifQ.Xd7vOBrbMuoxEwZpxWDH4bVVvPIFj60fzlt8Kd9Krwk")
	r.ServeHTTP(w, req)

	assert.Equal(t, "Authorization header format must be Bearer {token}\n", w.Body.String())
}
