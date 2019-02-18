package bundler

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

//Bundle maps file contents to keys, that can be used as strings in a Go
// application. File contents are stored as strings. Keys are defined as:
//[prefix]/path-from-[directory]/filename
//If [prefix] is empty string, the leading '/' is discarded
//For example,
// Given directory=file-bundler/test-files
// And prefix=TEST
// Then the key for "bacon.json" is TEST/bacon.json
//
// Given directory=file-bundler/test-files
// And prefix=""
// Then the key for "bacon.json" is bacon.json
//
//Prefixes are applied after the matcher and mapping
//
//matcher is a regular expression that determines which files are included. Each
// filename under "directory" is tested, and is included if matched. The matcher
// is tested against absolute filepaths. The matcher is applied before mapping.
//For example,
// Given directory=file-bundler/test-files
// And matcher=*.json
// Then "bacon.json" is bundled, and "usage.txt" is not
//
//mapping explicitly overrides file keying.
//For example,
// Given directory=file-bundler
// And mapping has the element: "test-files/bacon.json": "bacon"
// Then the key for "test-files/bacon.json" is "bacon"
//
//Default arguments:
//matcher=`.*`
//
//Note:
//While Bundle could be used for any file loading, it is intended to be used
// with `go generate`, and is designed to be seamless with `viper`
func Bundle(directory, matcher, prefix string, saveAsPlainText, compress bool,
	mapping map[string]string) (map[string]string, error) {
	bundle := make(map[string]string)
	rootDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}

	compiledMatcher, err := regexp.Compile(matcher)
	if err != nil {
		return nil, err
	}

	if err := filepath.Walk(rootDir, bundleWalkFn(rootDir, compiledMatcher,
		saveAsPlainText, compress, bundle)); err != nil {
		return nil, err
	}

	remap(bundle, mapping)

	if prefix != "" {
		prefixMapping := createPrefixedKeyRemapping(prefix, bundle)
		remap(bundle, prefixMapping)
	}

	return bundle, nil
}

//rootDir should be absolute path
func bundleWalkFn(rootDir string, matcher *regexp.Regexp, saveAsPlainText,
	compress bool, bundle map[string]string,
) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		if !matcher.MatchString(relPath) {
			return nil
		}

		rawFile, err := readFile(path)
		if err != nil {
			return err
		}

		if compress {
			rawFile = compressData(rawFile)
		}

		if saveAsPlainText {
			bundle[relPath] = string(rawFile)
		} else {
			bundle[relPath] = base64.StdEncoding.EncodeToString(rawFile)
		}
		return nil
	}
}

func compressData(data []byte) []byte {
	var b bytes.Buffer
	w, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
	_, _ = w.Write(data)
	_ = w.Close()
	return b.Bytes()
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	all, err := ioutil.ReadAll(f)
	return all, err
}

//remappings are defined as oldKey:newKey
//new keys must be unique wrt old keys
func remap(src, mapping map[string]string) {
	for oldKey, newKey := range mapping {
		src[newKey] = src[oldKey]
		delete(src, oldKey)
	}
}

//remappings are defined as oldKey:newKey
//creates a remapping map with prefixed keys
func createPrefixedKeyRemapping(prefix string,
	src map[string]string) map[string]string {
	p := make(map[string]string)

	for oldKey, _ := range src {
		newKey := filepath.Join(prefix, oldKey)
		p[oldKey] = newKey
	}

	return p
}

type bundleData struct {
	Package string
	Name    string
	Value   string
}

func WriteBundleWithViper(w io.Writer, goPackage, varName string,
	bundle map[string]string) error {
	return writeBundleWithTemplate(w, goPackage, varName, bundle,
		bundleWithViperTemplate)
}

func WriteBundle(w io.Writer, goPackage, varName string,
	bundle map[string]string) error {
	return writeBundleWithTemplate(w, goPackage, varName, bundle,
		bundleTemplate)
}

func writeBundleWithTemplate(w io.Writer, goPackage, varName string,
	bundle map[string]string, bundleTpl string) error {
	data := bundleData{
		Package: goPackage,
		Name:    varName,
		Value:   fmt.Sprintf("%#v", bundle),
	}
	tpl, err := template.New("").Parse(bundleTpl)
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
}
