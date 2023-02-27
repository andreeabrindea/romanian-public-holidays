package handler

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetResponseBody(t *testing.T) {
	_, err := GetResponseBody(-4, "abc")
	if err == nil {
		t.Errorf("expected an error, but got: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Hello, client")
		if err != nil {
			return
		}
	}))
	defer ts.Close()

	got, err := GetResponseBody(2024, ts.URL)
	if err != nil && got == nil {
		t.Errorf("unexpected error: %v", err)
	}
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
