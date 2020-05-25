package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_config_no_env_OK(t *testing.T) {
	ReadConfiguration()
	assert.Equal(t, true, CurConfiguration.Auth.AuthDisable == "true" || CurConfiguration.Auth.AuthDisable == "false")
}
func Test_config_overridden_OK(t *testing.T) {
	os.Setenv("auth.authdisable", "Overridden")
	ReadConfiguration()
	assert.Equal(t, "Overridden", CurConfiguration.Auth.AuthDisable)
}
