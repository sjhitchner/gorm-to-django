package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sjhitchner/gorm-to-django/django"
	"github.com/sjhitchner/gorm-to-django/gorm"

	"github.com/spf13/cobra"
)

const ()

func init() {
	rootCmd.Flags().StringP("gorm-models-dir", "g", "", "Path to GORM models")
	rootCmd.Flags().BoolP("django-admin", "a", false, "Generate Django Admin")
	rootCmd.Flags().BoolP("django-models", "m", true, "Generate Django Models")
	rootCmd.Flags().StringP("output-dir", "o", "", "Path to output Django models")
	rootCmd.Flags().StringP("template-dir", "t", "", "Path to template directory if you wish to override")
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

		outputDir, err := cmd.Flags().GetString("output-dir")
		if err != nil {
			return err
		}

		templateDir, err := cmd.Flags().GetString("template-dir")
		if err != nil {
			return err
		}

		generateModels, err := cmd.Flags().GetBool("django-models")
		if err != nil {
			return err
		}

		generateAdmin, err := cmd.Flags().GetBool("django-admin")
		if err != nil {
			return err
		}

		gormCh, err := gorm.Parse(gormDir)
		if err != nil {
			return err
		}

		var generator *django.Generator
		if templateDir == "" {
			generator, err = django.New()
			if err != nil {
				return err
			}
		} else {
			generator, err = django.NewWithCustomTemplates(templateDir)
			if err != nil {
				return err
			}
		}

		if err := generator.Build(gormCh); err != nil {
			return err
		}

		if generateModels {
			filename := filepath.Join(outputDir, "models.py")
			f, err := os.Create(filename)
			if err != nil {
				return err
			}

			if err := generator.GenerateModels(f); err != nil {
				return err
			}
		}

		if generateAdmin {
			filename := filepath.Join(outputDir, "admin.py")
			f, err := os.Create(filename)
			if err != nil {
				return err
			}

			if err := generator.GenerateAdmin(f); err != nil {
				return err
			}
		}

		return nil
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
