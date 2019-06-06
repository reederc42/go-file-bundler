package bundler

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFileDir = "../testdata"

func TestBundledKeys(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	if err != nil {
		t.Error(err)
	}

	bundle, err := Bundle(directory, "", "", true, false, nil)
	assert.NoError(t, err)
	assert.Contains(t, bundle, "bacon.json")
	assert.Contains(t, bundle, "usage.txt")
}

func TestMatchedKeys(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	if err != nil {
		t.Error(err)
	}

	bundle, err := Bundle(directory, ".*\\.json$", "", true, false, nil)
	assert.NoError(t, err)
	assert.Contains(t, bundle, "bacon.json")
	assert.NotContains(t, bundle, "usage.txt")
}

func TestMappedKeys(t *testing.T) {
	mapping := map[string]string{
		"bacon.json": "bacon",
	}

	directory, err := filepath.Abs(testFileDir)
	assert.NoError(t, err)

	bundle, err := Bundle(directory, "", "", true, false, mapping)
	assert.NoError(t, err)
	assert.Contains(t, bundle, "bacon")
	assert.Contains(t, bundle, "usage.txt")
}

func TestPrefixedKeys(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	assert.NoError(t, err)

	bundle, err := Bundle(directory, "", "FILE", true, false, nil)
	assert.NoError(t, err)

	assert.Contains(t, bundle, "FILE/bacon.json")
	assert.Contains(t, bundle, "FILE/usage.txt")
	assert.NotContains(t, bundle, "bacon.json")
	assert.NotContains(t, bundle, "usage.txt")
}

func TestContents(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	assert.NoError(t, err)

	baconJsonPath := filepath.Join(directory, "bacon.json")
	baconJsonContents, err := readFile(baconJsonPath)
	assert.NoError(t, err)

	usageTxtPath := filepath.Join(directory, "usage.txt")
	usageTxtContents, err := readFile(usageTxtPath)
	assert.NoError(t, err)

	bundle, err := Bundle(directory, "", "", true, false, nil)
	assert.NoError(t, err)

	assert.Contains(t, bundle, "bacon.json")
	bundledBacon, _ := bundle["bacon.json"]
	assert.Equal(t, string(baconJsonContents), bundledBacon)

	assert.Contains(t, bundle, "usage.txt")
	bundledUsage, _ := bundle["usage.txt"]
	assert.Equal(t, string(usageTxtContents), bundledUsage)
}

func TestRemap(t *testing.T) {
	src := map[string]string{
		"k1": "v1",
	}

	mapping := map[string]string{
		"k1": "key1",
	}

	remap(src, mapping)

	assert.Contains(t, src, "key1")
	assert.NotContains(t, src, "k1")
}
