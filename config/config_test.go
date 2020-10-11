package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigWrongPath(t *testing.T) {
	err := Read("./xxx.yaml")
	assert.NotNil(t, err)
}
func TestReadConfigWrongFile(t *testing.T) {
	err := Read("./data_test.go")
	assert.NotNil(t, err)
}

func TestReadConfig(t *testing.T) {
	err := Read("./config.yaml")
	assert.NotNil(t, Get())
	assert.Nil(t, err)

	assert.NotEmpty(t, GetLastModified())
	assert.NotEmpty(t, GetLastModifiedHash())
	assert.Equal(t, Get().Server.BasicAuth.UserName, "username")
	assert.Equal(t, Get().Server.BasicAuth.Password, "password")

	assert.NotNil(t, Get().CloudFoundries["cf-eu10"])
	assert.Equal(t, Get().CloudFoundries["cf-eu10"].UserName, "admin-eu10")
}
