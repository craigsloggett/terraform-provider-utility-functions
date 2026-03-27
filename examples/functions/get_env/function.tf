output "example" {
  description = "The value of a given environment variable."
  value       = provider::utilities::get_env("GOPATH")
}
