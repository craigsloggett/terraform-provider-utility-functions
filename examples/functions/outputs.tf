output "example" {
  description = "The output of a provider-defined function."
  value       = provider::utilities::generate_name("standard_name")
}
