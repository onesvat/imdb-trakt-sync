package trakt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cecobask/imdb-trakt-sync/internal/config"
)

func TestRefreshTokenUsesOAuthEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/oauth/token" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var body authTokensBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.GrantType != "refresh_token" || body.RefreshToken != "old-refresh" ||
			body.ClientID != "test-client" || body.ClientSecret != "test-secret" ||
			body.RedirectURI != "urn:ietf:wg:oauth:2.0:oob" || body.Code != "" {
			t.Errorf("unexpected refresh request body: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(authTokensResponse{
			AccessToken: "new-access", RefreshToken: "new-refresh", CreatedAt: 1700000000, ExpiresIn: 604800,
		}); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	clientID, clientSecret := "test-client", "test-secret"
	ac := newAuthClient(config.Trakt{ClientID: &clientID, ClientSecret: &clientSecret}, http.DefaultTransport)
	ac.baseURL = server.URL
	token, err := ac.getAccessToken(context.Background(), grantTypeRefreshToken, "old-refresh")
	if err != nil {
		t.Fatal(err)
	}
	if token.token != "new-access" || token.refreshToken != "new-refresh" || token.expiresAt.Unix() != 1700604800 {
		t.Fatalf("unexpected refreshed token: %+v", token)
	}
}
