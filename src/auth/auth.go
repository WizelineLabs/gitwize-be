package auth

import (
	"encoding/json"
	"fmt"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"gitwize-be/src/utils"
	"net/http"
	"os"
	"strings"
)

// JWTVerifier verifies if a token is valid
type JWTVerifier interface {
	Verify(token string, r *http.Request) bool
}

// OktaJWTVerifier verifies Okta tokens
type OktaJWTVerifier struct{}

// Verify verifies access token using Okta API
func (o OktaJWTVerifier) Verify(token string, r *http.Request) bool {
	tv := map[string]string{}
	tv["aud"] = "api://default"
	tv["cid"] = os.Getenv("SPA_CLIENT_ID")
	issuer := os.Getenv("ISSUER")
	jv := verifier.JwtVerifier{
		Issuer:           issuer,
		ClaimsToValidate: tv,
	}

	jvVerifier := jv.New()
	if _, err := jvVerifier.VerifyAccessToken(token); err != nil {
		return false
	}
	metaDataUrl := jvVerifier.Issuer + jvVerifier.Discovery.GetWellKnownUrl()

	resp, err := http.Get(metaDataUrl)

	if err != nil {
		fmt.Printf(utils.GetFuncName()+": Request for metadata was not successful: %s\n", err.Error())
		return false
	}

	defer resp.Body.Close()

	metaData := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&metaData)

	if respData, err := jvVerifier.Adaptor.Decode(token, metaData["jwks_uri"].(string)); err == nil {
		r.Header.Set("AuthenticatedUser", respData.(map[string]interface{})["sub"].(string))
		return true
	}
	return false
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
	return v.Verify(tokenParts[1], r)
}
