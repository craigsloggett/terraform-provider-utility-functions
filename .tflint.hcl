# Settings are described here: https://github.com/terraform-linters/tflint/blob/main/docs/user-guide/config.md
config {
  format              = "default"
  module              = true
  force               = false
  disabled_by_default = false
}

# Default Terraform ruleset described here: https://github.com/terraform-linters/tflint-ruleset-terraform/blob/main/docs/rules/README.md

# Disallow specifying a git repository as a module source without pinning to a version.
rule "terraform_module_pinned_source" {
  enabled = true
  style   = "semver"
}

# Checks that Terraform modules sourced from a registry specify a version.
rule "terraform_module_version" {
  enabled = true
  exact   = false
}

# Enforces naming conventions for resources, data sources, etc.
rule "terraform_naming_convention" {
  enabled = true
  format  = "snake_case"
}

plugin "terraform" {
  enabled = true
  version = "0.5.0"
  source  = "github.com/terraform-linters/tflint-ruleset-terraform"
  preset  = "all"
}
