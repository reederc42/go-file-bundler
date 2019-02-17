package bundler

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

const testFileDir = "../test-files"

func TestBundledKeys(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	if err != nil {
		t.Error(err)
	}

	bundle, err := Bundle(directory, "", "", true, false, nil)
	if err != nil {
		t.Error(err)
	}

	if _, ok := bundle["bacon.json"]; !ok {
		t.Errorf("did not find expected key: bacon.json")
	}

	if _, ok := bundle["usage.txt"]; !ok {
		t.Errorf("did not find expected key: usage.txt")
	}
}

func TestMatchedKeys(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	if err != nil {
		t.Error(err)
	}

	bundle, err := Bundle(directory, ".*\\.json$", "", true, false, nil)
	if err != nil {
		t.Error(err)
	}

	if _, ok := bundle["bacon.json"]; !ok {
		t.Errorf("did not find expected key: bacon.json")
	}

	if _, ok := bundle["usage.txt"]; ok {
		t.Errorf("found unexpected key: usage.txt")
	}
}

func TestMappedKeys(t *testing.T) {
	mapping := map[string]string{
		"bacon.json": "bacon",
	}

	directory, err := filepath.Abs(testFileDir)
	if err != nil {
		t.Error(err)
	}

	bundle, err := Bundle(directory, "", "", true, false, mapping)
	if err != nil {
		t.Error(err)
	}

	if _, ok := bundle["bacon"]; !ok {
		t.Errorf("did not find expected key: bacon")
	}

	if _, ok := bundle["usage.txt"]; !ok {
		t.Errorf("did not find expected key: usage.txt")
	}
}

func TestPrefixedKeys(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	if err != nil {
		t.Error(err)
	}

	bundle, err := Bundle(directory, "", "FILE", true, false, nil)
	if err != nil {
		t.Error(err)
	}

	if _, ok := bundle["FILE/bacon.json"]; !ok {
		t.Errorf("did not find expected key: FILE/bacon.json")
	}

	if _, ok := bundle["FILE/usage.txt"]; !ok {
		t.Errorf("did not find expected key: FILE/usage.txt")
	}

	if _, ok := bundle["bacon.json"]; ok {
		t.Errorf("found unexpected key: bacon.json")
	}

	if _, ok := bundle["usage.txt"]; ok {
		t.Errorf("found unexpected key: usage.txt")
	}
}

func TestContents(t *testing.T) {
	directory, err := filepath.Abs(testFileDir)
	if err != nil {
		t.Error(err)
	}

	baconJsonPath := filepath.Join(directory, "bacon.json")
	baconJsonContents, err := readFile(baconJsonPath)
	if err != nil {
		t.Error(err)
	}

	usageTxtPath := filepath.Join(directory, "usage.txt")
	usageTxtContents, err := readFile(usageTxtPath)
	if err != nil {
		t.Error(err)
	}

	bundle, err := Bundle(directory, "", "", true, false, nil)
	if err != nil {
		t.Error(err)
	}

	bundledBacon, ok := bundle["bacon.json"]
	if !ok {
		t.Errorf("did not find expected key: bacon.json")
	}

	if string(baconJsonContents) != bundledBacon {
		t.Errorf("bundled bacon content does not equal expected content")
	}

	bundledUsage, ok := bundle["usage.txt"]
	if !ok {
		t.Errorf("did not find expected key: usage.txt")
	}

	if string(usageTxtContents) != bundledUsage {
		t.Errorf("bundled usage content does not equal expected content")
	}
}

func TestRemap(t *testing.T) {
	src := map[string]string{
		"k1": "v1",
	}

	mapping := map[string]string{
		"k1": "key1",
	}

	remap(src, mapping)

	if _, ok := src["key1"]; !ok {
		t.Errorf("did not find expected key: key1")
	}

	if _, ok := src["k1"]; ok {
		t.Errorf("found unexpected key: k1")
	}
}

func TestWriteMap(t *testing.T) {
	packageName := "test_files_test"
	mapName := "bundle"
	bundle := map[string]string{
		"test-files/bacon.json": "ewogICJ1cmkiOiAiaHR0cHM6Ly9iYWNvbmlwc3VtLmNvbS9hcGkiLAogICJwYXJhbWV0ZXJzIjogewogICAgInR5cGUiOiAibWVhdC1hbmQtZmlsbGVyIgogIH0sCiAgInNwZWNpYWwtY2hhcmFjdGVyIjogImAiCn0=",
	}
	expectedFile := `//Generated File
package test_files_test

var bundle = map[string]string{
	"test-files/bacon.json": "ewogICJ1cmkiOiAiaHR0cHM6Ly9iYWNvbmlwc3VtLmNvbS9hcGkiLAogICJwYXJhbWV0ZXJzIjogewogICAgInR5cGUiOiAibWVhdC1hbmQtZmlsbGVyIgogIH0sCiAgInNwZWNpYWwtY2hhcmFjdGVyIjogImAiCn0=",
}
`
	var b bytes.Buffer
	err := WriteMap(&b, bundle, packageName, mapName)
	assert.NoError(t, err)
	assert.Equal(t, expectedFile, b.String())
}
