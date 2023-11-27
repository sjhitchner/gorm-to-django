package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const ()

func init() {
	rootCmd.Flags().StringP("gorm-models-dir", "g", "", "Path to GORM models")
	rootCmd.Flags().StringP("output-dir", "o", "", "Path to output Django models")
}

var rootCmd = &cobra.Command{
	Use:   "gorm-to-djanog",
	Short: "Generate Django models from GORM models",
	Long:  "Generate Django models from GORM models",
	RunE: func(cmd *cobra.Command, args []string) error {

		gormDir, err := cmd.Flags().GetString("gorm-models-dir")
		if err != nil {
			return err
		}

		/*
			outputDir, err := cmd.Flags().GetString("output-dir")
			if err != nil {
				return err
			}
		*/

		return parse(gormDir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
