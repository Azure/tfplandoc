/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
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
	Example: `tfplandoc generate /path/to/terraform/plan.json`,
	Args:    cobra.ExactArgs(1),
	Run:     runGenerateCmd,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	generateCmd.Flags().BoolP("all", "a", false, "Generate output for all resources, even with no changes")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runGenerateCmd(cmd *cobra.Command, args []string) {
	//planFile := cmd.Flag("planfile").Value.String()
	all := cmd.Flag("all").Value
	_ = all

	var inputReader io.Reader = cmd.InOrStdin()
	if len(args) > 0 && args[0] != "-" {
		file, err := os.Open(args[0])
		cobra.CheckErr(err)
		defer file.Close()
		inputReader = file
	}

	plan, err := readPlan(inputReader)
	cobra.CheckErr(err)

	resourceTable, err := generateResourceChangeTable(plan)
	cobra.CheckErr(err)

	err = md.NewMarkdown(os.Stdout).H1("Terraform Plan Documentation").LF().
		Table(resourceTable).
		LF().
		PlainText(fmt.Sprintf("Plan file: %s", path.Clean(args[0]))).
		LF().Build()
	cobra.CheckErr(err)
}

func readPlan(r io.Reader) (*terraform.PlanStruct, error) {
	planBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading plan: %w", err)
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
