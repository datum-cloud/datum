package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// PrivateKeyJwtKeyFileContent defines the structure of the JSON key file for private key JWT auth.
// This typically contains a user/service ID, a key ID, and the private key.
// Example fields: userId, keyId, key (PEM-encoded RSA private key)
type PrivateKeyJwtKeyFileContent struct {
	UserID string `json:"userId"` // The issuer and subject of the JWT assertion.
	KeyID  string `json:"keyId"`  // The key identifier, used in the JWT header's 'kid' claim.
	Key    string `json:"key"`    // The PEM-encoded private key (e.g., RSA) used for signing the JWT.
}

// PrivateKeyJwtTokenSource implements oauth2.TokenSource for service accounts
// using private key JWT authentication (RFC 7523).
// It generates a JWT assertion, signs it with a private key, and exchanges it for an OAuth2 token.
type PrivateKeyJwtTokenSource struct {
	logger logr.Logger

	keyFilePath string   // Path to the service account key file.
	tokenURL    string   // The URL of the OAuth2 token endpoint.
	audience    string   // The audience claim for the JWT assertion (usually the token endpoint or a resource server).
	scopes      []string // OAuth2 scopes to request.
	httpClient  *http.Client

	// Parsed key file content
	keyContent *PrivateKeyJwtKeyFileContent
	privateKey *rsa.PrivateKey

	// Mutex to protect token generation if multiple goroutines call Token() concurrently
	// on a non-cached source (though typically used with ReuseTokenSource).
	mu sync.Mutex
}

// NewPrivateKeyJwtTokenSource creates a new PrivateKeyJwtTokenSource.
// It reads and parses the key file upon creation to ensure it's valid.
func NewPrivateKeyJwtTokenSource(
	logger logr.Logger,
	keyFilePath string,
	tokenURL string,
	audience string,
	scopes []string,
	httpClient *http.Client,
) (*PrivateKeyJwtTokenSource, error) { // Return type changed back
	actLog := logger
	if actLog == (logr.Logger{}) {
		actLog = logr.Discard()
	}

	// Use http.DefaultClient if none provided
	actualHTTPClient := httpClient
	if actualHTTPClient == nil {
		actLog.V(1).Info("No HTTP client provided for token source, using http.DefaultClient.")
		actualHTTPClient = http.DefaultClient // Reverted: Use DefaultClient directly
	} else {
		actLog.V(1).Info("Using provided HTTP client for token source.")
	}

	content, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account key file %q: %w", keyFilePath, err)
	}

	var kfc PrivateKeyJwtKeyFileContent
	if err := json.Unmarshal(content, &kfc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service account key file %q: %w", keyFilePath, err)
	}

	if kfc.UserID == "" || kfc.KeyID == "" || kfc.Key == "" {
		return nil, fmt.Errorf("invalid service account key file %q: missing userId, keyId, or key", keyFilePath)
	}

	block, _ := pem.Decode([]byte(kfc.Key))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from private key in %q", keyFilePath)
	}

	if block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		actLog.Info("Found PEM block with potentially unsupported type for direct RSA parsing, attempting generic PKCS#8 then PKCS#1", "type", block.Type)
	}

	var parsedKey *rsa.PrivateKey
	if block.Type == "RSA PRIVATE KEY" {
		parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			actLog.V(1).Info("Failed to parse PEM block as PKCS#1 RSA private key, will attempt PKCS#8", "error", err.Error())
		}
	}

	if parsedKey == nil {
		pkcs8Key, errPkcs8 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if errPkcs8 != nil {
			originalRsaErrorMsg := "(no PKCS#1 attempt or error)"
			if err != nil { // err is from the x509.ParsePKCS1PrivateKey attempt
				originalRsaErrorMsg = err.Error()
			}
			return nil, fmt.Errorf("failed to parse private key from %q as PKCS#1 RSA ('%s') or PKCS#8 ('%s')", keyFilePath, originalRsaErrorMsg, errPkcs8.Error())
		}
		var ok bool
		parsedKey, ok = pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("parsed PKCS#8 key from %q is not an RSA private key (type: %T)", keyFilePath, pkcs8Key)
		}
	}

	return &PrivateKeyJwtTokenSource{
		logger:      actLog.WithName("PrivateKeyJwtTokenSource"),
		keyFilePath: keyFilePath,
		tokenURL:    tokenURL,
		audience:    audience,
		scopes:      scopes,
		httpClient:  actualHTTPClient, // Use the potentially wrapped client
		keyContent:  &kfc,
		privateKey:  parsedKey,
	}, nil
}

// Token retrieves a new token from the OAuth2 token endpoint using a private key JWT assertion.
func (s *PrivateKeyJwtTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info("Attempting to fetch new token using private key JWT")

	// 1. Create JWT assertion
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": s.keyContent.UserID,
		"sub": s.keyContent.UserID,
		"aud": s.audience,
		"iat": jwt.NewNumericDate(now),
		"exp": jwt.NewNumericDate(now.Add(5 * time.Minute)), // Short-lived assertion (e.g., 5 minutes)
	}

	assertionToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	assertionToken.Header["kid"] = s.keyContent.KeyID

	signedAssertion, err := assertionToken.SignedString(s.privateKey)
	if err != nil {
		s.logger.Error(err, "Failed to sign JWT assertion")
		return nil, fmt.Errorf("failed to sign JWT assertion: %w", err)
	}
	s.logger.V(1).Info("Successfully signed JWT assertion")

	// 2. Make POST request to s.tokenURL
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", signedAssertion)
	if len(s.scopes) > 0 {
		form.Set("scope", strings.Join(s.scopes, " ")) // Ensure scope is added!
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Add a timeout to the HTTP request
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", s.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		s.logger.Error(err, "Failed to create token request")
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use the client stored in the struct
	s.logger.V(1).Info("Sending token request to OAuth2 token endpoint", "url", s.tokenURL)
	resp, err := s.httpClient.Do(req) // httpClient is assigned during New... from actualHTTPClient
	if err != nil {
		s.logger.Error(err, "Failed to send token request to OAuth2 token endpoint")
		return nil, fmt.Errorf("failed to send token request to OAuth2 token endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error(fmt.Errorf("unexpected status code: %s", resp.Status), "OAuth2 token request failed")
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			s.logger.Error(readErr, "Failed to read error response body from token endpoint")
			return nil, fmt.Errorf("oauth2 token request failed with status %s and unable to read response body: %w", resp.Status, readErr)
		}
		return nil, fmt.Errorf("oauth2 token request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	// 3. Parse response and return oauth2.Token
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		s.logger.Error(err, "Failed to decode token response")
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		s.logger.Error(nil, "Token response missing access_token")
		return nil, fmt.Errorf("token response missing access_token")
	}

	s.logger.Info("Successfully fetched token from OAuth2 token endpoint", "type", tokenResp.TokenType, "expires_in_seconds", tokenResp.ExpiresIn)

	return &oauth2.Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		Expiry:      time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}, nil
}
