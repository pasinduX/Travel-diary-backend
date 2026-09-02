package integrations

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth0JWKS struct {
	Keys []Auth0JWK `json:"keys"`
}

type Auth0JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var (
	auth0JWKSCacheMu sync.Mutex
	auth0JWKSCache   = map[string]cachedJWKS{}
)

type cachedJWKS struct {
	keys    map[string]*rsa.PublicKey
	fetched time.Time
}

func ParseAuth0Token(ctx context.Context, jwksURL, issuer, audience, tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}
		keys, err := fetchAuth0Keys(ctx, jwksURL)
		if err != nil {
			return nil, err
		}
		key, ok := keys[kid]
		if !ok {
			return nil, fmt.Errorf("unknown kid: %s", kid)
		}
		return key, nil
	}, jwt.WithIssuer(issuer), jwt.WithAudience(audience))
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, errors.New("invalid auth0 token")
	}
	return token, claims, nil
}

func fetchAuth0Keys(ctx context.Context, jwksURL string) (map[string]*rsa.PublicKey, error) {
	auth0JWKSCacheMu.Lock()
	if cached, ok := auth0JWKSCache[jwksURL]; ok && time.Since(cached.fetched) < 10*time.Minute {
		keys := cached.keys
		auth0JWKSCacheMu.Unlock()
		return keys, nil
	}
	auth0JWKSCacheMu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks request failed: %s", resp.Status)
	}

	var jwks Auth0JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	keys := make(map[string]*rsa.PublicKey, len(jwks.Keys))
	for _, jwk := range jwks.Keys {
		if jwk.Kid == "" || jwk.N == "" || jwk.E == "" {
			continue
		}
		pub, err := jwkToPublicKey(jwk.N, jwk.E)
		if err != nil {
			continue
		}
		keys[jwk.Kid] = pub
	}

	auth0JWKSCacheMu.Lock()
	auth0JWKSCache[jwksURL] = cachedJWKS{keys: keys, fetched: time.Now().UTC()}
	auth0JWKSCacheMu.Unlock()

	return keys, nil
}

func jwkToPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}
	exponent := 0
	for _, b := range eBytes {
		exponent = exponent<<8 + int(b)
	}
	if exponent == 0 {
		return nil, errors.New("invalid exponent")
	}
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: exponent,
	}, nil
}
