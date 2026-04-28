package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: len(cafeList["moscow"])},
	}

	for _, req := range requests {
		response := httptest.NewRecorder()

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?city=moscow&count=%d", req.count), nil)

		handler.ServeHTTP(response, request)

		require.Equal(t, http.StatusOK, response.Code)

		body := response.Body.String()

		var cafes []string
		if body == "" {
			cafes = []string{}
		} else {
			cafes = strings.Split(body, ",")
		}

		assert.Equalf(t, req.want, len(cafes), "count=%d", req.count)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, req := range requests {
		response := httptest.NewRecorder()

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/cafe?city=moscow&search=%s", req.search), nil)

		handler.ServeHTTP(response, request)

		require.Equal(t, http.StatusOK, response.Code)

		body := response.Body.String()

		var cafes []string
		if body == "" {
			cafes = []string{}
		} else {
			cafes = strings.Split(body, ",")
		}

		assert.Equalf(t, req.wantCount, len(cafes), "search=%q", req.search)

		for _, cafe := range cafes {
			cafe = strings.ToLower(cafe)
			search := strings.ToLower(req.search)

			assert.True(t, strings.Contains(cafe, search))
		}
	}
}
