package plex

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebsocketReadErrorLogLevel(t *testing.T) {
	tests := []struct {
		name      string
		closeCode int
		wantLevel string
	}{
		{name: "normal close", closeCode: websocket.CloseNormalClosure, wantLevel: `"level":"DEBUG"`},
		{name: "going away", closeCode: websocket.CloseGoingAway, wantLevel: `"level":"DEBUG"`},
		{name: "unexpected close", closeCode: websocket.CloseInternalServerErr, wantLevel: `"level":"ERROR"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
				if err != nil {
					t.Errorf("upgrade websocket: %v", err)
					return
				}
				defer conn.Close()

				if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(tt.closeCode, "")); err != nil {
					t.Errorf("send close frame: %v", err)
				}
			}))
			defer server.Close()

			var output bytes.Buffer
			client := &Plex{URL: server.URL}
			client.SetLogger(slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})))

			var callbackErr error
			client.SubscribeToNotifications(NewNotificationEvents(), make(chan os.Signal), func(err error) {
				callbackErr = err
			})

			if !websocket.IsCloseError(callbackErr, tt.closeCode) {
				t.Fatalf("callback error = %v, want websocket close code %d", callbackErr, tt.closeCode)
			}
			if !strings.Contains(output.String(), tt.wantLevel) {
				t.Errorf("log = %q, want %s", output.String(), tt.wantLevel)
			}
		})
	}
}
