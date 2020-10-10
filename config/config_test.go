package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigWrongPath(t *testing.T) {
	cfs, err := New("./xxx.yaml")
	assert.Nil(t, cfs)
	assert.NotNil(t, err)
}
func TestReadConfigWrongFile(t *testing.T) {
	cfs, err := New("./data_test.go")
	assert.Nil(t, cfs)
	assert.NotNil(t, err)
}

func TestReadConfig(t *testing.T) {
	cfs, err := New("./config.yaml")
	assert.NotNil(t, cfs)
	assert.Nil(t, err)
	assert.NotEqual(t, uint32(0), cfs.GetLastModifiedHash())
	assert.NotEqual(t, time.Time{}, cfs.GetLastModified())
}

func TestHasChanged(t *testing.T) {
	cfs, err := New("./config.yaml")
	assert.NotNil(t, cfs)
	assert.Nil(t, err)
	unchanged := hasChanged("./config.yaml")
	assert.False(t, unchanged)
	changed := hasChanged("./data_test.go")
	assert.True(t, changed)
}

func TestDeepCopy(t *testing.T) {
	cfg1, err := New("./config.yaml")
	assert.NotNil(t, cfg1)
	assert.Nil(t, err)
	cfg2, err := New("./config.yaml")
	assert.NotNil(t, cfg2)
	assert.Nil(t, err)

	assert.Equal(t, cfg1, cfg2)

	cfg2.Server.BasicAuth.UserName = "bla"
	assert.NotEqual(t, cfg1, cfg2)
}

func TestContent(t *testing.T) {
	cfg, err := New("./config.yaml")
	assert.NotNil(t, cfg)
	assert.Nil(t, err)

	assert.Equal(t, cfg.Server.BasicAuth.UserName, "username")
	assert.Equal(t, cfg.Server.BasicAuth.Passowrd, "password")

	assert.NotNil(t, cfg.CloudFoundries["cf-eu10"])
	assert.Equal(t, cfg.CloudFoundries["cf-eu10"].UserName, "admin-eu10")
}
