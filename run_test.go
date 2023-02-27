package main

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
func TestConvertTOJSONWithDate(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		expected []time.Time
	}{
		{
			name: "valid response",
			data: []byte(`
				{
					"response": {
						"holidays": [
							{
								"date": {
									"iso": "2023-03-01"
								}
							},
							{
								"date": {
									"iso": "2023-04-01"
								}
							}
						]
					}
				}`),
			expected: []time.Time{
				time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "invalid json",
			data: []byte(`
				{
					"response": {
				}`),
			expected: []time.Time{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := ConvertTOJSONWithDate(tc.data)
			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
