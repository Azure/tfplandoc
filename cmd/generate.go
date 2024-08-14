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
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	md "github.com/nao1215/markdown"
	"github.com/spf13/cobra"
)

var generateAllFlag, textFlag bool

var generateSymbols = map[tfjson.Action]rune{
	tfjson.ActionCreate:     '🟢',
	tfjson.ActionRead:       '🔵',
	tfjson.ActionUpdate:     '🟠',
	tfjson.ActionDelete:     '🔴',
	tfjson.ActionNoop:       '⚪',
	tfjson.Action("forget"): '🟣',
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate documentation for Terraform plan output",
	Long: `Generates documentation from the Terraform plan output. Has resource/output create, read, update, delete and forget operations.

Key:

🟢 - Create
🔵 - Read
🟠 - Update
🔴 - Delete
⚪ - Noop
🟣 - Forget (remove form state but do not destroy)

`,
	Example: `tfplandoc generate /path/to/terraform/plan.json`,
	Args:    cobra.ExactArgs(1),
	Run:     runGenerateCmd,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&generateAllFlag, "all", "a", false, "Generate output for all resources, even with no changes")
	generateCmd.Flags().BoolVarP(&textFlag, "text", "t", false, "Output in plain text, without symbols and color")

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

	output := md.NewMarkdown(os.Stdout)
	output = output.H4("Resource Changes").LF()
	addResourceChangeTable(output, plan)
	output = output.H4("Output Changes").LF()
	addOutputChangeTable(output, plan)
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
	err = json.Unmarshal(planBytes, plan)
	cobra.CheckErr(err)

	return plan, nil
}

func addResourceChangeTable(m *md.Markdown, plan *tfjson.Plan) {
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
		if !generateAllFlag && rc.Change.Actions[0] == tfjson.ActionNoop {
			continue
		}
		actionStr := changeActionsToSymbolString(rc.Change.Actions)
		changeRows = append(changeRows, []string{rc.Address, actionStr})
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

func addOutputChangeTable(m *md.Markdown, plan *tfjson.Plan) {
	keys := make([]string, 0, len(plan.OutputChanges))
	for k := range plan.OutputChanges {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	changeRows := make([][]string, 0, len(plan.OutputChanges))
	for _, key := range keys {
		actions := plan.OutputChanges[key].Actions
		if len(actions) == 0 {
			continue
		}
		if !generateAllFlag && actions[0] == tfjson.ActionNoop {
			continue
		}
		actionStr := changeActionsToSymbolString(actions)
		changeRows = append(changeRows, []string{key, actionStr})
	}
	if len(changeRows) == 0 {
		m.PlainText("***No changes detected***").LF()
		return
	}
	var table md.TableSet
	table.Rows = changeRows
	table.Header = []string{"Output", "Change"}
	m.Table(table)
}

func changeActionsToSymbolString(actions []tfjson.Action) string {
	if textFlag {
		actionsStr := make([]string, 0, len(actions))
		for _, action := range actions {
			actionsStr = append(actionsStr, string(action))
		}
		return strings.Join(actionsStr, ",")
	}
	var res string
	for _, action := range actions {
		res += string(generateSymbols[action])
	}
	return res
}
