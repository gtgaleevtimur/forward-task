package handler

import (
	"forward/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	db := repository.NewRepository()
	router := NewRouter(db)
	assert.NotNil(t, router)
}

func TestController_Get(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		db := repository.NewRepository()
		router := NewRouter(db)
		ts := httptest.NewServer(router)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/check/101010101", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})
	t.Run("Negative with another method ", func(t *testing.T) {
		db := repository.NewRepository()
		router := NewRouter(db)
		ts := httptest.NewServer(router)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/check/101010101", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "method does not allowed", string(body))
	})
	t.Run("Positive with another route", func(t *testing.T) {
		db := repository.NewRepository()
		router := NewRouter(db)
		ts := httptest.NewServer(router)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/v1/", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, "route does not exist", string(body))
	})
}
