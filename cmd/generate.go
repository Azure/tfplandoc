/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"slices"

	md "github.com/go-spectest/markdown"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/spf13/cobra"
)

var generateAllFlag bool

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
	generateCmd.Flags().BoolVarP(&generateAllFlag, "all", "a", false, "Generate output for all resources, even with no changes")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runGenerateCmd(cmd *cobra.Command, args []string) {
	var inputReader io.Reader = cmd.InOrStdin()
	if len(args) > 0 && args[0] != "-" {
		file, err := os.Open(args[0])
		cobra.CheckErr(err)
		defer file.Close()
		inputReader = file
	}

	plan, err := readPlan(inputReader)
	cobra.CheckErr(err)

	output := md.NewMarkdown(os.Stdout).H2("Terraform Plan Documentation").LF()
	addResourceChangeTable(output, plan, false, generateAllFlag)
	output.LF().
		PlainText(fmt.Sprintf("Plan file: %s", path.Clean(args[0]))).
		LF().Build()
	cobra.CheckErr(err)
}

func readPlan(r io.Reader) (*tfjson.Plan, error) {
	planBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading plan: %w", err)
	}
	plan := new(tfjson.Plan)
	json.Unmarshal(planBytes, plan)

	if err != nil {
		return nil, fmt.Errorf("error parsing plan file %w", err)
	}
	return plan, nil
}

func addResourceChangeTable(m *md.Markdown, plan *tfjson.Plan, markdown, all bool) {
	slices.SortFunc(plan.ResourceChanges, func(a, b *tfjson.ResourceChange) int {
		// negative if a < b
		if a.Address < b.Address {
			return -1
		}
		// positive if a > b
		if a.Address > b.Address {
			return 1
		}
		// zero if a == b
		return 0
	})
	changeRows := make([][]string, 0, len(plan.ResourceChanges))
	for _, rc := range plan.ResourceChanges {
		if len(rc.Change.Actions) == 0 {
			continue
		}
		if !all && rc.Change.Actions[0] == tfjson.ActionNoop {
			continue
		}
		changeRows = append(changeRows, []string{rc.Address, fmt.Sprintf("%s", rc.Change.Actions)})
	}
	if len(changeRows) == 0 {
		m.PlainText("***No changes detected***").LF()
		return
	}
	var table md.TableSet
	table.Rows = changeRows
	table.Header = []string{"Resource", "Change"}
	m.Table(table)
}
