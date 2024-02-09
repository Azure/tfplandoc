/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"

	md "github.com/go-spectest/markdown"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate documentation for Terraform plan output",
	Long:    `Generates documentation from the Terraform plan output. Has resource creations, modifications and deletions.`,
	Example: `tfplandoc generate -p /path/to/terraform/plan.json`,
	Run:     runGenerateCmd,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateCmd.PersistentFlags().StringP("planfile", "p", "", "Path to terraform plan file")
	generateCmd.MarkFlagFilename("planfile", "json")
	generateCmd.MarkFlagRequired("planfile")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runGenerateCmd(cmd *cobra.Command, args []string) {
	planFile := cmd.Flag("planfile").Value.String()
	if planFile == "" {
		cmd.Help()
		cobra.CheckErr(errors.New("plan file is required"))
	}

	plan, err := readPlanFile(planFile)
	cobra.CheckErr(err)

	resourceTable, err := generateResourceChangeTable(plan)
	cobra.CheckErr(err)

	err = md.NewMarkdown(os.Stdout).H1("Terraform Plan Documentation").LF().
		Table(resourceTable).
		LF().
		PlainText(fmt.Sprintf("Plan file: %s", path.Clean(planFile))).
		LF().Build()
	cobra.CheckErr(err)
}

func readPlanFile(planFile string) (*terraform.PlanStruct, error) {
	planBytes, err := os.ReadFile(planFile)
	if err != nil {
		return nil, fmt.Errorf("error reading plan file %w", err)
	}

	plan, err := terraform.ParsePlanJSON(string(planBytes))
	if err != nil {
		return nil, fmt.Errorf("error parsing plan file %w", err)
	}
	return plan, nil
}

func generateResourceChangeTable(plan *terraform.PlanStruct) (md.TableSet, error) {
	var result md.TableSet
	changeRows := make([][]string, len(plan.ResourceChangesMap))
	resourceMapKeys := maps.Keys(plan.ResourceChangesMap)
	sort.Strings(resourceMapKeys)
	i := 0
	for _, key := range resourceMapKeys {
		changeRows[i] = []string{key, fmt.Sprintf("%s", plan.ResourceChangesMap[key].Change.Actions)}
		i++
	}
	result.Header = []string{"Resource", "Change"}
	result.Rows = changeRows
	return result, nil
}
