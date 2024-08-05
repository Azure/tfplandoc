# tfplandoc - a tool to generate Terraform plan documentation

This tool generates documentation for Terraform plans. It reads the plan file and generates a markdown file with the resources that will be created, updated, or deleted.

## Installation

Download the compiled release from the [releases page](https://github.com/matt-FFFFFF/tfplandoc/releases), or install it using `go install`:

```bash
go install github.com/matt-FFFFFF/tfplandoc@latest
```

## Usage

First generate a plan file using `terraform plan -out <file>`, then convert to JSON using `terraform show -json`.
Finally, run `tfplandoc` with the plan file as an argument:

```bash
terraform plan -out tfplan && terraform show -json tfplan >tfplan.json
```

```bash
tfplandoc generate tfplan.json
```

You can also pipe the output of `terraform show -json` directly to `tfplandoc`:

```bash
terraform plan -out tfplan && terraform show -json tfplan | tfplandoc generate -
```

### Output key

To show the key for the output changes, use the `--help` or `-h` flag:

```text
Generates documentation from the Terraform plan output. Has resource creations, modifications and deletions.

Key:

🟢 - Create
🔵 - Read
🟠 - Update
🔴 - Delete
⚪ - Noop
🟣 - Forget (remove form state but do not destroy)

Usage:
  tfplandoc generate [flags]

Examples:
tfplandoc generate /path/to/terraform/plan.json

Flags:
  -a, --all    Generate output for all resources, even with no changes
  -h, --help   help for generate
  -t, --text   Output in plain text, without symbols and color
```


### Showing all output

By default, the tool will only show resources and outputs that have changes.
To show all resources and outputs, use the `--all` flag:

```bash
tfplandoc generate tfplan.json --all
```

## Example output

```text
## Terraform Plan Documentation

### Resource Changes

|                                                    RESOURCE                                                     | CHANGE |
|-----------------------------------------------------------------------------------------------------------------|--------|
| module.management.azapi_resource.data_collection_rule["change_tracking"]                                        | 🟢     |
| module.management.azapi_resource.data_collection_rule["defender_sql"]                                           | 🟢     |
| module.management.azapi_resource.data_collection_rule["vm_insights"]                                            | 🟢     |
| module.management.azapi_resource.sentinel_onboarding[0]                                                         | 🟢     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/AgentHealthAssessment"]       | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/AntiMalware"]                 | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/ChangeTracking"]              | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/SQLAdvancedThreatProtection"] | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/SQLAssessment"]               | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/SQLVulnerabilityAssessment"]  | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/Security"]                    | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/ServiceMap"]                  | 🔴     |
| module.management.azurerm_log_analytics_solution.management["Microsoft/OMSGallery/Updates"]                     | 🔴     |
| module.management.azurerm_log_analytics_solution.security_insights_for_removal                                  | 🟣     |
| module.management.azurerm_user_assigned_identity.management["ama"]                                              | 🟢     |

### Output Changes

|            OUTPUT             | CHANGE |
|-------------------------------|--------|
| test_data_collection_rule_ids | 🟢     |
| test_managed_identity_ids     | 🟢     |


Plan file: tfplan.json
```
