package route

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func TestSetPixelsRejectsInvalidGatewayUser(t *testing.T) {
	app := New(&gorm.DB{})

	for _, userID := range []string{"", "0", "-1", "invalid", "999999999999999999999"} {
		t.Run(userID, func(t *testing.T) {
			body := strings.NewReader(`{"pixels":[{"x":0,"y":0,"color":"ffffff"}]}`)
			request := httptest.NewRequest(http.MethodPost, "/api/pixels", body)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Neko-User-Id", "42")
			if userID != "" {
				request.Header.Set("X-NekoBox-User-ID", userID)
			}
			response := httptest.NewRecorder()
			app.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusUnauthorized, response.Body)
			}
		})
	}
}
