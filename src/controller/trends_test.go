package controller

import (
	"github.com/stretchr/testify/assert"
	"gitwize-be/src/configuration"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func Test_QuarterlyTrends_RejectedPR(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"

	from := time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix()
	to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.Local).Unix()
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+"/trends?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10),
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)

	// must contain data of June
	//resp := QuarterlyTrends{
	//	PercentageRejectedPR: make(map[string]string),
	//}
	//json.Unmarshal(w.Body.Bytes(), &resp)
	//assert.NotNil(t, resp.PercentageRejectedPR["June"])
}
