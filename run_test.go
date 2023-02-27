package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetResponseBody(t *testing.T) {
	t.Run("returns error for unsupported year", func(t *testing.T) {
		year := -1
		_, err := getResponseBody(year)
		expectedErr := errors.New("unsupported year")
		if err.Error() != expectedErr.Error() {
			t.Errorf("got %v, expected %v", err, expectedErr)
		}
	})

	t.Run("returns response body for supported year", func(t *testing.T) {
		year := 2024
		expectedResponseBody := []byte("test response body")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(expectedResponseBody)
			if err != nil {
				return
			}
		})
		server := httptest.NewServer(handler)
		defer server.Close()
		_, err := getResponseBody(year)
		if err != nil {
			t.Errorf("got error %v, expected nil", err)
		}
	})
}
