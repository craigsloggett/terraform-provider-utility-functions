# Utility Functions

The Utility Functions provider is used to add utility functions to a Terraform project. The provider has no configuration.

## Required Secrets

The following secrets must be added to the GitHub repository:

| Secret Name                  | Description                                                                                     |
| ---------------------------- | ----------------------------------------------------------------------------------------------- |
| `GPG_PRIVATE_KEY`            | The GPG private key used to sign provider releases before publishing to the Terraform registry. |
| `GPG_PRIVATE_KEY_PASSPHRASE` | The passphrase for the GPG private signing key.                                                 |
