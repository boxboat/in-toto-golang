package main

import (
	"fmt"
	"os"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/spf13/cobra"
)

var (
	spiffeUDS  string
	layoutPath string
	keyPath    string
	certPath   string
	key        intoto.Key
	cert       intoto.Key
)

var rootCmd = &cobra.Command{
	Use:   "in-toto",
	Short: "Framework to secure integrity of software supply chains",
	Long:  `A framework to secure the integrity of software supply chains https://in-toto.io/`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&spiffeUDS,
		"spiffe-workload-api-path",
		"",
		"uds path for spiffe workload api",
	)
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
