package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	intoto "github.com/boxboat/in-toto-golang/in_toto"
	"github.com/spf13/cobra"
)

var stepName string
var keyPath string
var certPath string
var materialsPaths []string
var productsPaths []string
var outDir string
var lStripPaths []string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executes the passed command and records paths and hashes of 'materials'",
	Long: `Executes the passed command and records paths and hashes of 'materials' (i.e.
files before command execution) and 'products' (i.e. files after command
execution) and stores them together with other information (executed command,
return value, stdout, stderr, ...) to a link metadata file, which is signed
with the passed key.  Returns nonzero value on failure and zero otherwise.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		//Load Key
		var cert, key intoto.Key

		if spiffeUDS != "" {
			ctx := context.Background()
			key = intoto.GetSVID(ctx, spiffeUDS)

		} else {

			if err := key.LoadKeyDefaults(keyPath); err != nil {
				fmt.Println("Invalid Key Error:", err.Error())
				os.Exit(1)
			}

			if len(certPath) > 0 {
				if err := cert.LoadKeyDefaults(certPath); err != nil {
					fmt.Println("Invalid Certificate Error:", err.Error())
					os.Exit(1)
				}

				key.KeyVal.Certificate = cert.KeyVal.Certificate
			}
		}

		block, err := intoto.InTotoRun(stepName, materialsPaths, productsPaths, args, key, []string{"sha256"}, []string{}, lStripPaths)
		if err != nil {
			fmt.Println("Error generating meta-block:", err.Error())
			os.Exit(1)
		}

		linkName := fmt.Sprintf(intoto.LinkNameFormat, block.Signed.(intoto.Link).Name, key.KeyID)

		err = block.Dump(filepath.Join(outDir, linkName))
		if err != nil {
			fmt.Println("Error writing meta-block:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&stepName,
		"name", "n", "",
		`Name used to associate the resulting link metadata
with the corresponding step defined in an in-toto
layout.`)
	runCmd.Flags().StringVarP(&keyPath,
		"key", "k", "",
		`Path to a PEM formatted private key file used to sign
the resulting link metadata. (passing one of '--key'
or '--gpg' is required) `)
	runCmd.Flags().StringArrayVarP(&materialsPaths,
		"materials", "m", []string{},
		`Paths to files or directories, whose paths and hashes
are stored in the resulting link metadata before the
command is executed. Symlinks are followed.`)
	runCmd.Flags().StringArrayVarP(&productsPaths,
		"products", "p", []string{},
		`Paths to files or directories, whose paths and hashes
are stored in the resulting link metadata after the
command is executed. Symlinks are followed.`)
	runCmd.Flags().StringVarP(&certPath,
		"cert", "c", "",
		`Path to a PEM formatted certificate that corresponds with
the provided key.`)
	runCmd.Flags().StringVarP(&outDir,
		"output-directory", "d", "./",
		`directory to store link metadata`)
	runCmd.Flags().StringArrayVar(&lStripPaths, "lstrip-paths", []string{},
		`path prefixes used to left-strip artifact paths before storing
them to the resulting link metadata. If multiple prefixes
are specified, only a single prefix can match the path of
any artifact and that is then left-stripped. All prefixes
are checked to ensure none of them are a left substring
of another.`)

	runCmd.MarkFlagRequired("name")
}
