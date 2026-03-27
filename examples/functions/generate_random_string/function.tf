output "example" {
  description = "The result of generating a random string of length 10."
  value       = provider::utilities::generate_random_string(10)
}
