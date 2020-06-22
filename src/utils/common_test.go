package utils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"regexp"
	"testing"
	"time"
)

func Test_TimeTrack(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	expectedResult := ".*gitwize-be/src/utils.Test_TimeTrack took .*s"

	TimeTrack(time.Now(), GetFuncName())
	assert.Regexp(t, regexp.MustCompile(expectedResult), buf.String())
}
func Test_GetFuncName(t *testing.T) {
	nameFunction := "gitwize-be/src/utils.Test_GetFuncName"
	assert.Equal(t, nameFunction, GetFuncName())
}

func Test_Trace(t *testing.T) {
	expectedResult := "Entering: .*:.*:\\d+"
	assert.Regexp(t, regexp.MustCompile(expectedResult), Trace())
}
