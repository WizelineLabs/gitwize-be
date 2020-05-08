package auth

import (
	"net/http"
	"os"
	"strings"
	verifier "github.com/okta/okta-jwt-verifier-golang"
)

// JWTVerifier verifies if a token is valid
type JWTVerifier interface {
	Verify(token string) bool
}

// OktaJWTVerifier verifies Okta tokens
type OktaJWTVerifier struct {}

// Verify verifies access token using Okta API
func (o OktaJWTVerifier) Verify(token string) bool {
	tv := map[string]string{}
	tv["aud"] = "api://default"
	tv["cid"] = os.Getenv("SPA_CLIENT_ID")
	issuer := os.Getenv("ISSUER")
	jv := verifier.JwtVerifier{
		Issuer:           issuer,
		ClaimsToValidate: tv,
	}

	_, err := jv.New().VerifyAccessToken(token)

	if err != nil {
		return false
	}

	return true
}

// IsAuthorized verifies access token passed in Authorization header
func IsAuthorized(v JWTVerifier, r *http.Request) bool {
	if v == nil {
		v = OktaJWTVerifier{}
	}
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}

	tokenParts := strings.Split(authHeader, "Bearer ")
	return v.Verify(tokenParts[1])
}
