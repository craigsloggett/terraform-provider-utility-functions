output "get_env" {
  # value     = provider::utilities::get_env("GOPATH")
  description = "The value of a given environment variable."
  value       = "Uncomment output value to test provider function." # TFLint doesn't yet support provider-defined functions.
}

output "generate_random_string" {
  # value     = provider::utilities::generate_random_string(10)
  description = "A random number with a specified length."
  value       = "Uncomment output value to test provider function." # TFLint doesn't yet support provider-defined functions.
}
