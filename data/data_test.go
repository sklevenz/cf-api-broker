package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigWrongPath(t *testing.T) {
	cfs, err := NewCloudFoundryMetaData("./xxx.yaml")
	assert.Nil(t, cfs)
	assert.NotNil(t, err)
}
func TestReadConfigWrongFile(t *testing.T) {
	cfs, err := NewCloudFoundryMetaData("./data_test.go")
	assert.Nil(t, cfs)
	assert.NotNil(t, err)
}

func TestReadConfig(t *testing.T) {
	cfs, err := NewCloudFoundryMetaData("./config.yaml")
	assert.NotNil(t, cfs)
	assert.Nil(t, err)
}

func TestHasChanged(t *testing.T) {
	cfs, err := NewCloudFoundryMetaData("./config.yaml")
	assert.NotNil(t, cfs)
	assert.Nil(t, err)
	unchanged := hasChanged("./config.yaml")
	assert.False(t, unchanged)
	changed := hasChanged("./data_test.go")
	assert.True(t, changed)
}

func TestDeepCopy(t *testing.T) {
	cfs1, err := NewCloudFoundryMetaData("./config.yaml")
	assert.NotNil(t, cfs1)
	assert.Nil(t, err)
	cfs2, err := NewCloudFoundryMetaData("./config.yaml")
	assert.NotNil(t, cfs2)
	assert.Nil(t, err)

	assert.Equal(t, cfs1, cfs2)
	cfs2.CloudFoundries[0].Password = "***********"
	cfs2.CloudFoundries[0].Labels = []string{"232"}
	assert.NotEqual(t, cfs1, cfs2)
}
