package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigWrongPath(t *testing.T) {
	cfs, err := NewConfig("./xxx.yaml")
	assert.Nil(t, cfs)
	assert.NotNil(t, err)
}
func TestReadConfigWrongFile(t *testing.T) {
	cfs, err := NewConfig("./data_test.go")
	assert.Nil(t, cfs)
	assert.NotNil(t, err)
}

func TestReadConfig(t *testing.T) {
	cfs, err := NewConfig("./config.yaml")
	assert.NotNil(t, cfs)
	assert.Nil(t, err)
	assert.NotEqual(t, uint32(0), cfs.GetLastModifiedHash())
	assert.NotEqual(t, time.Time{}, cfs.GetLastModified())
}

func TestHasChanged(t *testing.T) {
	cfs, err := NewConfig("./config.yaml")
	assert.NotNil(t, cfs)
	assert.Nil(t, err)
	unchanged := hasChanged("./config.yaml")
	assert.False(t, unchanged)
	changed := hasChanged("./data_test.go")
	assert.True(t, changed)
}

func TestDeepCopy(t *testing.T) {
	cfg1, err := NewConfig("./config.yaml")
	assert.NotNil(t, cfg1)
	assert.Nil(t, err)
	cfg2, err := NewConfig("./config.yaml")
	assert.NotNil(t, cfg2)
	assert.Nil(t, err)

	assert.Equal(t, cfg1, cfg2)

	cfg2.Broker.BasicAuth.UserName = "bla"
	assert.NotEqual(t, cfg1, cfg2)
}
