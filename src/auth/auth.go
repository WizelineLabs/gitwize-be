package auth

import (
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"gitwize-be/src/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

// JWTVerifier verifies if a token is valid
type JWTVerifier interface {
	Verify(token string, r *http.Request) (bool, string)
}

// OktaJWTVerifier verifies Okta tokens
type OktaJWTVerifier struct{}

// Verify verifies access token using Okta API
func (o OktaJWTVerifier) Verify(token string, r *http.Request) (bool, string) {
	tv := map[string]string{}
	tv["aud"] = "api://default"
	tv["cid"] = os.Getenv("SPA_CLIENT_ID")
	issuer := os.Getenv("ISSUER")
	jv := verifier.JwtVerifier{
		Issuer:           issuer,
		ClaimsToValidate: tv,
	}

	if jwt, err := jv.New().VerifyAccessToken(token); err == nil {
		userEmail := jwt.Claims["sub"].(string)
		r.Header.Set("AuthenticatedUser", userEmail)
		log.Println(utils.GetFuncName() + ": AuthenticatedUser=" + userEmail)
		return true, userEmail
	}
	return false, ""
}

// IsAuthorized verifies access token passed in Authorization header
func IsAuthorized(v JWTVerifier, r *http.Request) (bool, string) {
	if v == nil {
		v = OktaJWTVerifier{}
	}
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false, ""
	}

	tokenParts := strings.Split(authHeader, "Bearer ")
	return v.Verify(tokenParts[1], r)
}
