package main

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestSession(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login2", nil)
	router.ServeHTTP(w, req)
	req0, _ := http.NewRequest("GET", "/login2", nil)
	req0.Header.Add("Cookie", strings.Join(w.HeaderMap["Set-Cookie"],";"))

	router.ServeHTTP(w, req0)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"count\":0}{\"count\":1}", w.Body.String())
}