//go:build integration

package auth

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/oauth2"
)

func TestPrivateKeyJwtTokenSource_Integration(t *testing.T) {
	// --- Configuration from Environment Variables ---
	keyFile := os.Getenv("TEST_PKJWT_KEY_FILE")
	if keyFile == "" {
		t.Skip("Skipping integration test: TEST_PKJWT_KEY_FILE environment variable not set.")
		return
	}

	tokenURL := os.Getenv("TEST_PKJWT_TOKEN_URL")
	if tokenURL == "" {
		tokenURL = "https://auth.staging.env.datum.net/oauth/v2/token" // Default to staging
		t.Logf("Using default TEST_PKJWT_TOKEN_URL: %s", tokenURL)
	}

	audience := os.Getenv("TEST_PKJWT_AUDIENCE")
	if audience == "" {
		audience = "https://auth.staging.env.datum.net" // Default to staging
		t.Logf("Using default TEST_PKJWT_AUDIENCE: %s", audience)
	}

	scopesStr := os.Getenv("TEST_PKJWT_SCOPES")
	if scopesStr == "" {
		scopesStr = "openid,email" // Default scopes
		t.Logf("Using default TEST_PKJWT_SCOPES: %s", scopesStr)
	}
	scopes := []string{}
	if len(strings.TrimSpace(scopesStr)) > 0 {
		scopes = strings.Split(scopesStr, ",")
	}
	// --- End Configuration ---

	logger := logr.Discard() // Or use testr.New(t) if more detailed test logging is needed

	// Create the token source
	t.Logf("Creating token source with key file: %s, tokenURL: %s, audience: %s, scopes: %v",
		keyFile, tokenURL, audience, scopes)
	tokSrc, err := NewPrivateKeyJwtTokenSource(logger, keyFile, tokenURL, audience, scopes, nil /* use default http client with logging */)
	if err != nil {
		t.Fatalf("NewPrivateKeyJwtTokenSource() error = %v", err)
	}

	// Wrap with caching (optional but good practice, mirrors serve.go)
	cachedSrc := oauth2.ReuseTokenSource(nil, tokSrc)

	// Attempt to fetch a token
	t.Log("Attempting to fetch initial token...")
	token, err := cachedSrc.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	// Basic validation
	if token == nil {
		t.Fatal("Token() returned nil token")
	}

	if token.AccessToken == "" {
		t.Error("Token() returned token with empty AccessToken")
	}

	if token.Expiry.IsZero() {
		t.Error("Token() returned token with zero Expiry")
	} else if !token.Valid() {
		// Check if expiry is reasonable (e.g., not already expired)
		t.Errorf("Token() returned token that is already expired or invalid. Expiry: %s", token.Expiry.Format(time.RFC3339))
	} else {
		t.Logf("Successfully retrieved token. Type: %s, Expires: %s", token.TokenType, token.Expiry.Format(time.RFC3339))
	}

	// Optional: Try fetching again to test caching/reuse
	t.Log("Attempting to fetch token again (testing reuse)...")
	token2, err := cachedSrc.Token()
	if err != nil {
		t.Fatalf("Second Token() call error = %v", err)
	}
	if token2 == nil {
		t.Fatal("Second Token() call returned nil token")
	}
	if token.AccessToken != token2.AccessToken {
		// This might happen if the first token was very close to expiry, not necessarily an error
		t.Logf("Second token has different AccessToken, potentially refreshed.")
	} else {
		t.Logf("Successfully reused token.")
	}
}
