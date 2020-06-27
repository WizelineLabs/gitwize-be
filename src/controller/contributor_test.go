package controller

import (
	"github.com/stretchr/testify/assert"
	"gitwize-be/src/configuration"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// Test direct curl http://0.0.0.0:8000/api/v1/repositories/1/contributor?date_from=1277788580&date_to=1593148580&author_email=test@wizeline.com
func TestGetContributorStats(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"

	to := time.Now().Unix()
	from := to - 10*365*24*3600
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+"/contributor?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10)+"&author_email=test@wizeline.com", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	log.Println("query ", gwEndPointRepository+
		strconv.Itoa(1)+"/contributor?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10)+"&author_email=test@wizeline.com")

	expectedResult := "{\"repository_id\":\\d+," +
		"\"author_email\":\".*\"," +
		"\"author_name\":\".*\"," +
		"\"commits\":\\d+," +
		"\"addition_loc\":\\d+," +
		"\"deletion_loc\":\\d+," +
		"\"num_files\":\\d+," +
		"\"loc_percent\":\\d+," +
		"\"date\":\".*\"}"
	assert.Regexp(t, regexp.MustCompile(expectedResult), w.Body.String())
}
