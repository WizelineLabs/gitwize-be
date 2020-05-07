package main

import (
	"testing"
	"net/http"
	"gitwize-be/src/auth"
)

// --- auth.go tests ---
type mockJWTVerifier struct{}

func (m mockJWTVerifier) Verify(token string) bool {
	if token == "valid-token" {
		return true
	}
	return false
}

func Test_Request_Authorized(t *testing.T) {
	r := &http.Request {
		Header: map[string][]string{
					"Authorization": {"Bearer valid-token"},
		},
	}
	v := mockJWTVerifier{}
	authorized := auth.IsAuthorized(v, r)
	if !authorized {
		t.Error("Expected authorized!")
	}
}

func Test_Request_UnAuthorized(t *testing.T) {
	r := &http.Request {
		Header: map[string][]string{
					"Authorization": {"Bearer XXXXXX"},
		},
	}
	v := mockJWTVerifier{}
	authorized := auth.IsAuthorized(v, r)
	if authorized {
		t.Error("Expected unauthorized!")
	}
}
// --- end auth.go tests ---