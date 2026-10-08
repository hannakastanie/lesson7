package main

import (
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
		{"/cafe?city=tula&count=-1", http.StatusBadRequest, "incorrect count"},
		{"/cafe?city=tula&count=%206", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, v.status, response.Code)
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

		require.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	count := len(cafeList["moscow"])

	requests := []struct {
		request   string
		wantCount int
	}{
		{"/cafe?count=0&city=moscow", 0},
		{"/cafe?count=1&city=moscow", 1},
		{"/cafe?count=2&city=moscow", 2},
		{"/cafe?count=100&city=moscow", min(count, 100)},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		answer := response.Body.String()

		require.Equal(t, http.StatusOK, response.Code)
		if answer == "" {
			assert.Len(t, answer, v.wantCount)
		} else {
			result := strings.Split(answer, ",")
			assert.Len(t, result, v.wantCount)
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request    string
		wantCount  int
		wantSearch string
	}{
		{"/cafe?city=moscow&search=фасоль", 0, ""},
		{"/cafe?city=moscow&search=кофе", 2, "кофе"},
		{"/cafe?city=moscow&search=вилка", 1, "вилка"},
		{"/cafe?city=moscow&search=ФЛЮГЕГЕХАЙМЕН", 0, ""},
		{"/cafe?city=moscow&count=1&search=кофе", 1, "кофе"},
		{"/cafe?city=moscow&search=кофе&count=1", 1, "кофе"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		answer := response.Body.String()

		require.Equal(t, http.StatusOK, response.Code)
		if answer == "" {
			assert.Len(t, answer, v.wantCount)
		} else {
			result := strings.Split(answer, ",")
			assert.Len(t, result, v.wantCount)
			for _, eachResult := range result {
				assert.True(t, strings.Contains(
					strings.ToLower(eachResult),
					strings.ToLower(v.wantSearch),
				))
			}
		}
	}
}
