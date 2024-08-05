package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	md "github.com/nao1215/markdown"
	"github.com/stretchr/testify/assert"
)

func TestChangeActionsToSymbolString(t *testing.T) {
	// Test case 1: Single action
	actions1 := []tfjson.Action{tfjson.ActionCreate}
	expected1 := "🟢"
	result1 := changeActionsToSymbolString(actions1)
	if result1 != expected1 {
		t.Errorf("Expected %s, but got %s", expected1, result1)
	}

	// Test case 2: Multiple actions
	actions2 := []tfjson.Action{tfjson.ActionCreate, tfjson.ActionUpdate, tfjson.ActionDelete}
	expected2 := "🟢🟠🔴"
	result2 := changeActionsToSymbolString(actions2)
	if result2 != expected2 {
		t.Errorf("Expected %s, but got %s", expected2, result2)
	}

	// Test case 3: No actions
	actions3 := []tfjson.Action{}
	expected3 := ""
	result3 := changeActionsToSymbolString(actions3)
	if result3 != expected3 {
		t.Errorf("Expected %s, but got %s", expected3, result3)
	}

	// Test case 4: textFlag is true
	textFlag = true
	expected4 := "create,update,delete"
	result4 := changeActionsToSymbolString(actions2)
	if result4 != expected4 {
		t.Errorf("Expected %s, but got %s", expected4, result4)
	}
}

func TestRunGenerateCmd(t *testing.T) {

	// 	expected := `## Terraform Plan Documentation

	// ### Resource Changes

	// |       RESOURCE       | CHANGE |
	// |----------------------|--------|
	// | terraform_data.one   | 🟢     |
	// | terraform_data.three | 🟢     |
	// | terraform_data.two   | 🟢     |

	// ### Output Changes

	// ***No changes detected***

	// Plan file: ../testdata/basic/tfplan.json`

	args := []string{"../testdata/basic/tfplan.json"}
	runGenerateCmd(generateCmd, args)
}

func TestReadPlan(t *testing.T) {
	// Test case 1: Valid plan
	planBytes1 := []byte(`{"format_version":"1.2","terraform_version":"1.7.2","planned_values":{"root_module":{"resources":[{"address":"terraform_data.one","mode":"managed","type":"terraform_data","name":"one","provider_name":"terraform.io/builtin/terraform","schema_version":0,"values":{"id":"ad3f762d-d8d1-b3bf-3d87-c14a36e75859","input":"one","output":"one","triggers_replace":"one"},"sensitive_values":{}},{"address":"terraform_data.three","mode":"managed","type":"terraform_data","name":"three","provider_name":"terraform.io/builtin/terraform","schema_version":0,"values":{"id":"e0f84f20-51ef-27f4-9f40-43a693f33e39","input":"three-modified","output":"three-modified","triggers_replace":null},"sensitive_values":{}}]}},"resource_changes":[{"address":"terraform_data.one","mode":"managed","type":"terraform_data","name":"one","provider_name":"terraform.io/builtin/terraform","change":{"actions":["no-op"],"before":{"id":"ad3f762d-d8d1-b3bf-3d87-c14a36e75859","input":"one","output":"one","triggers_replace":"one"},"after":{"id":"ad3f762d-d8d1-b3bf-3d87-c14a36e75859","input":"one","output":"one","triggers_replace":"one"},"after_unknown":{},"before_sensitive":{},"after_sensitive":{}}},{"address":"terraform_data.three","mode":"managed","type":"terraform_data","name":"three","provider_name":"terraform.io/builtin/terraform","change":{"actions":["no-op"],"before":{"id":"e0f84f20-51ef-27f4-9f40-43a693f33e39","input":"three-modified","output":"three-modified","triggers_replace":null},"after":{"id":"e0f84f20-51ef-27f4-9f40-43a693f33e39","input":"three-modified","output":"three-modified","triggers_replace":null},"after_unknown":{},"before_sensitive":{},"after_sensitive":{}}}],"prior_state":{"format_version":"1.0","terraform_version":"1.7.2","values":{"root_module":{"resources":[{"address":"terraform_data.one","mode":"managed","type":"terraform_data","name":"one","provider_name":"terraform.io/builtin/terraform","schema_version":0,"values":{"id":"ad3f762d-d8d1-b3bf-3d87-c14a36e75859","input":"one","output":"one","triggers_replace":"one"},"sensitive_values":{}},{"address":"terraform_data.three","mode":"managed","type":"terraform_data","name":"three","provider_name":"terraform.io/builtin/terraform","schema_version":0,"values":{"id":"e0f84f20-51ef-27f4-9f40-43a693f33e39","input":"three-modified","output":"three-modified","triggers_replace":null},"sensitive_values":{}}]}}},"configuration":{"provider_config":{"terraform":{"name":"terraform","full_name":"terraform.io/builtin/terraform"}},"root_module":{"resources":[{"address":"terraform_data.one","mode":"managed","type":"terraform_data","name":"one","provider_config_key":"terraform","expressions":{"input":{"constant_value":"one"},"triggers_replace":{"constant_value":"one"}},"schema_version":0},{"address":"terraform_data.three","mode":"managed","type":"terraform_data","name":"three","provider_config_key":"terraform","expressions":{"input":{"constant_value":"three-modified"}},"schema_version":0}]}},"timestamp":"2024-02-16T14:53:28Z","errored":false}`)

	_, err1 := readPlan(bytes.NewReader(planBytes1))
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}
}

func TestAddResourceChangeTable(t *testing.T) {
	// Test case 1: No changes detected
	m1 := md.NewMarkdown(os.Stdout)
	plan1 := &tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{},
	}
	addResourceChangeTable(m1, plan1)
	expected1 := "***No changes detected***"
	result1 := m1.String()
	result1 = strings.TrimSpace(result1)
	assert.Equal(t, expected1, result1)

	// Test case 2: Changes detected
	m2 := md.NewMarkdown(os.Stdout)
	plan2 := &tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{
				Address: "terraform_data.one",
				Change: &tfjson.Change{
					Actions: []tfjson.Action{tfjson.ActionCreate},
				},
			},
			{
				Address: "terraform_data.two",
				Change: &tfjson.Change{
					Actions: []tfjson.Action{tfjson.ActionUpdate},
				},
			},
		},
	}
	addResourceChangeTable(m2, plan2)
	expected2 := "|      RESOURCE      | CHANGE |\n|--------------------|--------|\n| terraform_data.one | 🟢     |\n| terraform_data.two | 🟠     |"
	result2 := m2.String()
	result2 = strings.TrimSpace(result2)
	assert.Equal(t, expected2, result2)
}
