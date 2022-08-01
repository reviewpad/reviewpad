// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "reviewpad-cli",
	Long: "reviewpad-cli is command line interface to run reviewpad commands.",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&reviewpadFile, "file", "f", "", "input reviewpad file")
	rootCmd.MarkPersistentFlagRequired("file")
	rootCmd.SilenceUsage = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
