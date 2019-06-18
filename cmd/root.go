package cmd

import (
	"os"
	"path/filepath"

	bundler "github.com/reederc42/go-file-bundler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "file-bundler",
	Short: "bundles files into Go code",
	Long:  help,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		m, err := bundler.Bundle(viper.GetString(optionSrcDirectory),
			viper.GetString(optionMatcher), viper.GetString(optionPrefix),
			viper.GetBool(optionPlainText), viper.GetBool(optionCompress),
			viper.GetBool(optionHTTPPaths), nil)
		if err != nil && !viper.GetBool(optionSuppressErrors) {
			return err
		} else if err != nil && viper.GetBool(optionSuppressErrors) {
			m = make(map[string]string)
		}
		f, _ := os.Create(viper.GetString(optionDstFile))
		if viper.GetBool(optionViper) {
			err = bundler.WriteBundleWithViper(f,
				viper.GetString(optionPackage), viper.GetString(optionMapName),
				m)
		} else {
			err = bundler.WriteBundle(f, viper.GetString(optionPackage),
				viper.GetString(optionMapName), m)
		}
		return err
	},
}

func Execute() {
	_ = rootCmd.Execute()
}

func init() {
	defaultSrcDir, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	defaultOutputFile, err := filepath.Abs("bundle.go")
	if err != nil {
		panic(err)
	}
	rootCmd.Flags().StringP(optionSrcDirectory, "d", defaultSrcDir,
		"source directory")
	rootCmd.Flags().StringP(optionDstFile, "o", defaultOutputFile,
		"output file")
	rootCmd.Flags().StringP(optionMatcher, "m", ".*", "file matcher")
	rootCmd.Flags().StringP(optionPrefix, "x", "", "key prefix")
	rootCmd.Flags().StringP(optionPackage, "p", "", "package name")
	rootCmd.MarkFlagRequired(optionPackage)
	rootCmd.Flags().BoolP(optionPlainText, "t", false,
		"save as plain text instead of base64")
	rootCmd.Flags().StringP(optionMapName, "n",
		"bundle", "name of generated map")
	rootCmd.Flags().BoolP(optionCompress, "g", false,
		"use best gzip compression")
	rootCmd.Flags().BoolP(optionViper, "v", false, "integrate with viper")
	rootCmd.Flags().Bool(optionSuppressErrors, false,
		"on error writes empty output")
	rootCmd.Flags().Bool(optionHTTPPaths, false,
		"sets path separator to http /")
}

var help = `Bundle maps file contents to keys, that can be used as strings in a Go
	application. File contents are stored as strings. Keys are defined as:
[prefix]/path-from-[directory]/filename
If [prefix] is empty string, the leading '/' is discarded
For example,
	Given directory=file-bundler/test-files
	And prefix=TEST
	Then the key for "bacon.json" is TEST/bacon.json

	Given directory=file-bundler/test-files
	And prefix=""
	Then the key for "bacon.json" is bacon.json

Prefixes are applied after the matcher and mapping

matcher is a regular expression that determines which files are included. Each
	filename under "directory" is tested, and is included if matched. The matcher
	is tested against absolute filepaths. The matcher is applied before mapping.
For example,
	Given directory=file-bundler/test-files
	And matcher=*.json
	Then "bacon.json" is bundled, and "usage.txt" is not

mapping explicitly overrides file keying.
For example,
	Given directory=file-bundler
	And mapping has the element: "test-files/bacon.json": "bacon"
	Then the key for "test-files/bacon.json" is "bacon"

Default arguments:
matcher='.*'

Note:
While file-bundler could be used for any file loading, it is intended to be used
	with 'go generate' and 'viper'`
