package controller

import (
	"github.com/stretchr/testify/assert"
	"gitwize-be/src/configuration"
	"net/http"
	"regexp"
	"strconv"
	"testing"
)

// Test direct curl http://0.0.0.0:8000/api/v1/repositories/1/contributor?date_from=1277788580&date_to=1593148580&author_email=test@wizeline.com
func TestGetContributorStats(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"

	to := int64(2524637115) // 2050
	from := int64(0)        // 1970
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+"/contributor?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10),
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)
	expectedResult := "{\"repository_id\":\\d+," +
		"\"email\":\".*\"," +
		"\"name\":\".*\"," +
		"\"commits\":\\d+," +
		"\"additions\":\\d+," +
		"\"deletions\":\\d+," +
		"\"filesChange\":\\d+," +
		"\"changePercent\":\\d+," +
		"\"date\":\".*\"}"
	assert.Regexp(t, regexp.MustCompile(expectedResult), w.Body.String())
}
