package cmd

import (
	"github.com/reederc42/file-bundler/bundler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   "file-bundler",
	Short: "bundles files into Go code",
	Long:  "bundles files into Go code",
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		m, err := bundler.Bundle(viper.GetString(optionSrcDirectory),
			viper.GetString(optionMatcher), viper.GetString(optionPrefix),
			viper.GetBool(optionPlainText), viper.GetBool(optionCompress), nil)
		if err != nil {
			return err
		}
		f, _ := os.Create(viper.GetString(optionDstFile))
		err = bundler.WriteMap(f, m, viper.GetString(optionPackage),
			viper.GetString(optionMapName))
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
}
