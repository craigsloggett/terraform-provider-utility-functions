output "example" {
  # value     = provider::utilities::get_env("GOPATH")
  description = "The value of a given environment variable."
  value       = "Uncomment output value to test provider function." # TFLint doesn't yet support provider-defined functions.
}
